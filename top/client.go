// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package top

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gorilla/websocket"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	gops "github.com/mitchellh/go-ps"
)

// Client struct
type Client struct {
	options       *ClientOptions
	ui            terminal.UI
	cliConnection plugin.CliConnection
	eventrouting  *eventrouting.EventRouter
	router        *eventrouting.EventRouter
}

// ClientOptions needed to start the Client
type ClientOptions struct {
	Debug      bool
	NoTopCheck bool
	Cygwin     bool
	Nozzles    int
}

// NewClient instantiating the top client
func NewClient(cliConnection plugin.CliConnection, options *ClientOptions, ui terminal.UI) *Client {

	return &Client{
		options:       options,
		ui:            ui,
		cliConnection: cliConnection,
	}
}

// Ask Created my own Ask func instead of the CF provided one because the CF version adds
// a call to PromptColor(">") which does not work cleanly on MS-Windows
func (c *Client) Ask(prompt string) string {
	fmt.Printf("%s ", prompt)
	rd := bufio.NewReader(os.Stdin)
	line, err := rd.ReadString('\n')
	if err == nil {
		return strings.TrimSpace(line)
	}
	return ""
}

// Start starting the client
func (c *Client) Start() {

	if !c.options.NoTopCheck && c.shouldExitTop() {
		// There are other instances of top running and user requested to exit
		return
	}

	toplog.SetDebugEnabled(c.options.Debug)

	conn := c.cliConnection

	isLoggedIn, err := conn.IsLoggedIn()
	if err != nil {
		c.ui.Failed(err.Error())
		return
	}
	if !isLoggedIn {
		c.ui.Failed("Must login first")
		return
	}

	fmt.Printf("Loading...")
	fmt.Printf("\r           \r")

	ui := ui.NewMasterUI(conn)

	c.router = ui.GetRouter()

	toplog.Info("Top started at " + time.Now().Format("01-02-2006 15:04:05"))

	c.setupFirehoseConnections()

	ui.Start()
}

func (c *Client) setupFirehoseConnections() {
	privileged, err := c.hasDopplerFirehoseScope()
	if err != nil {
		c.ui.Failed("Could not determine privileges. Are you logged in?", err)
		return
	}

	if privileged {
		toplog.Info("Running with doppler.firehose privileges - opening %v nozzles", c.options.Nozzles)
		subscriptionID := "TopPlugin_" + util.Pseudo_uuid()
		for i := 0; i < c.options.Nozzles; i++ {
			go c.createAndKeepAliveNozzle(subscriptionID, "", i)
		}
		toplog.Info("Starting %v firehose nozzle instances", c.options.Nozzles)
		return
	}

	apps, err := c.cliConnection.GetApps()
	if err != nil {
		c.ui.Failed("Fetching all Apps failed: %v", err)
		return
	}
	toplog.Info("Running without doppler.firehose privileges - opening %v nozzles", len(apps))
	for i, application := range apps {
		go c.createAndKeepAliveNozzle("", application.Guid, i)
		toplog.Info("Starting app nozzle #%v instance for App %s", i, application.Name)
	}
}

func (c *Client) createAndKeepAliveNozzle(subscriptionID, appGUID string, instanceID int) error {
	minRetrySeconds := (2 * time.Second)

	for {
		// This is a blocking call if no error
		startTime := time.Now()
		var err error
		if len(subscriptionID) > 0 {
			err = c.createNozzle(subscriptionID, instanceID)
		} else {
			err = c.createAppNozzle(appGUID, instanceID)
		}
		if err != nil {
			errMsg := err.Error()
			notAuthorized := strings.Contains(errMsg, "authorized")
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) || notAuthorized {
				toplog.Error("Nozzle #%v - Stopped with error: %v", instanceID, err)
				if notAuthorized {
					toplog.Error("Are you sure you have 'admin' privileges on foundation?")
					toplog.Error("See needed permissions for this plugin here:")
					toplog.Error("https://github.com/ECSTeam/cloudfoundry-top-plugin")
				}
				break
			}
			toplog.Warn("Nozzle #%v - error: %v", instanceID, err)
		}
		toplog.Warn("Nozzle #%v - Shutdown. Nozzle instance will be restarted", instanceID)
		lastRetry := time.Now().Sub(startTime)
		if lastRetry < minRetrySeconds {
			toplog.Info("Nozzle #%v - Nozzle instance restart too fast, delaying for %v", instanceID, minRetrySeconds)
			time.Sleep(minRetrySeconds)
		}
	}
	return nil
}

func (c *Client) createNozzle(subscriptionID string, instanceID int) error {
	// Delay each nozzle instance creation by 1 second
	time.Sleep(time.Duration(instanceID) * time.Second)

	conn := c.cliConnection
	dopplerEndpoint, err := conn.DopplerEndpoint()
	if err != nil {
		return err
	}

	skipVerifySSL, err := conn.IsSSLDisabled()
	if err != nil {
		return err
	}

	dopplerConnection := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipVerifySSL}, nil)

	tokenRefresher := NewTokenRefresher(conn, instanceID)
	dopplerConnection.RefreshTokenFrom(tokenRefresher)
	dopplerConnection.SetMinRetryDelay(500 * time.Millisecond)
	dopplerConnection.SetMaxRetryDelay(15 * time.Second)
	dopplerConnection.SetIdleTimeout(15 * time.Second)

	authToken, err := conn.AccessToken()
	if err != nil {
		return err
	}

	messages, errors := dopplerConnection.Firehose(subscriptionID, authToken)
	defer dopplerConnection.Close()

	toplog.Info("Nozzle #%v - Started", instanceID)

	eventError := c.routeEvents(instanceID, messages, errors)
	if eventError != nil {
		msg := eventError.Error()
		if strings.Contains(msg, "Invalid authorization") {
			return eventError
		}
	}
	return nil
}

func (c *Client) createAppNozzle(appGUID string, instanceID int) error {

	// Delay each nozzle instance creation by 1 second
	time.Sleep(time.Duration(instanceID) * time.Second)

	conn := c.cliConnection
	dopplerEndpoint, err := conn.DopplerEndpoint()
	if err != nil {
		return err
	}

	skipVerifySSL, err := conn.IsSSLDisabled()
	if err != nil {
		return err
	}

	dopplerConnection := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipVerifySSL}, nil)

	tokenRefresher := NewTokenRefresher(conn, instanceID)
	dopplerConnection.RefreshTokenFrom(tokenRefresher)

	authToken, err := conn.AccessToken()
	if err != nil {
		return err
	}

	messages, errors := dopplerConnection.StreamWithoutReconnect(appGUID, authToken)
	defer dopplerConnection.Close()

	toplog.Info("Nozzle #%v for %s - Started", instanceID, appGUID)

	eventError := c.routeEvents(instanceID, messages, errors)
	if eventError != nil {
		msg := eventError.Error()
		if strings.Contains(msg, "Invalid authorization") {
			return eventError
		}
	}
	return nil
}

func (c *Client) routeEvents(instanceID int, messages <-chan *events.Envelope, errors <-chan error) error {
	for {
		select {
		case envelope := <-messages:
			c.router.Route(instanceID, envelope)
		case err := <-errors:
			c.handleError(instanceID, err)
			// Nozzle connection does not seem to recover from errors well, so
			// return here so it can be closed and a new instanced opened
			return err
		}
	}
}

func (c *Client) handleError(instanceID int, err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		msg := fmt.Sprintf("Nozzle #%v - Normal Websocket Closure: %v", instanceID, err)
		toplog.Error(msg)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		msg := fmt.Sprintf("Nozzle #%v - Disconnected because nozzle couldn't keep up (CloseError): %v", instanceID, err)
		toplog.Error(msg)
	default:
		msg := fmt.Sprintf("Nozzle #%v - Error reading firehose: %v", instanceID, err)
		toplog.Error(msg)
	}
}

func (c *Client) shouldExitTop() bool {
	numRunning := c.getNumberOfTopPluginsRunning() - 1
	if numRunning > 0 {
		plural := ""
		if numRunning > 1 {
			plural = "s"
		}
		fmt.Printf("Currently %v instance%v of cf top already running on this OS.\n", numRunning, plural)
		for {
			response := c.Ask("Do you want to continue and start top? (y/n)")
			response = strings.ToLower(response)
			if response == "y" || response == "yes" {
				break
			} else if response == "n" || response == "no" {
				fmt.Printf("Top will not be started\n")
				return true
			}
		}
	}
	return false
}

// Check how many instance of this plugin are running on the currently
// OS.  This is based on checking the process list for processes that
// have the same name as our plugin process name.
func (c *Client) getNumberOfTopPluginsRunning() int {

	p, err := gops.Processes()
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	if len(p) <= 0 {
		fmt.Printf("should have processes")
	}

	processName := ""
	for _, p1 := range p {
		if os.Getpid() == p1.Pid() {
			processName = p1.Executable()
		}
	}

	numberRunning := 0
	for _, p1 := range p {
		if p1.Executable() == processName {
			//fmt.Printf("Found: %v  PID: %v\n", p1.Executable(), p1.Pid())
			numberRunning++
		}
	}

	//fmt.Printf("Number of programs running with my same name: %v\n", numberRunning)
	return numberRunning
}

func (c *Client) hasDopplerFirehoseScope() (bool, error) {
	authToken, err := c.cliConnection.AccessToken()
	if err != nil {
		return false, err
	}

	decodedAccessToken, err := decodeAccessToken(authToken)
	if err != nil {
		return false, err
	}

	jsonParsed, err := gabs.ParseJSON(decodedAccessToken)
	if err != nil {
		return false, err
	}

	scopes, err := jsonParsed.Search("scope").Children()
	if err != nil {
		return false, err
	}

	for _, scope := range scopes {
		if scope.Data().(string) == "doppler.firehose" {
			return true, nil
		}
	}

	return false, nil
}

func decodeAccessToken(accessToken string) (tokenJSON []byte, err error) {
	tokenParts := strings.Split(accessToken, " ")

	if len(tokenParts) < 2 {
		return
	}

	token := tokenParts[1]
	encodedParts := strings.Split(token, ".")

	if len(encodedParts) < 3 {
		return
	}

	encodedTokenJSON := encodedParts[1]
	return base64Decode(encodedTokenJSON)
}

func base64Decode(encodedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(restorePadding(encodedData))
}

func restorePadding(seg string) string {
	switch len(seg) % 4 {
	case 2:
		seg = seg + "=="
	case 3:
		seg = seg + "="
	}
	return seg
}
