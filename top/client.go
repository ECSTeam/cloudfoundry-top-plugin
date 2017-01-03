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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gorilla/websocket"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventrouting"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui"
	gops "github.com/mitchellh/go-ps"
)

type Client struct {
	options       *ClientOptions
	ui            terminal.UI
	cliConnection plugin.CliConnection
	eventrouting  *eventrouting.EventRouter
	router        *eventrouting.EventRouter
}

type ClientOptions struct {
	AppGUID    string
	Debug      bool
	NoTopCheck bool
	Cygwin     bool
	Nozzles    int
}

func NewClient(cliConnection plugin.CliConnection, options *ClientOptions, ui terminal.UI) *Client {

	return &Client{
		options:       options,
		ui:            ui,
		cliConnection: cliConnection,
	}
}

// Created mine own Ask func instead of the CF provided one because the CF version adds
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
	// We request an access token to confirm that authentication has not expired
	_, err = conn.AccessToken()
	if err != nil {
		c.ui.Failed("AccessToken failed: %v", err)
		return
	}
	fmt.Printf("\r           \r")

	ui := ui.NewMasterUI(conn)

	c.router = ui.GetRouter()

	toplog.Info("Top started at " + time.Now().Format("01-02-2006 15:04:05"))

	apps, err := c.cliConnection.GetApps()
	if err != nil {
		c.ui.Failed("Fetching all Apps failed: %v", err)
		return
	}

	for i, application := range apps {
		c.createNozzles(application.Guid, application.Name, i)
	}

	ui.Start()

}

func (c *Client) createNozzles(appGUID, appName string, instanceID int) {
	go c.createAndKeepAliveNozzle(appGUID, instanceID)
	toplog.Info("Starting nozzle #%v instance for App %s", instanceID, appName)
}

func (c *Client) createAndKeepAliveNozzle(appGUID string, instanceID int) error {

	minRetrySeconds := (2 * time.Second)

	for {
		// This is a blocking call if no error
		startTime := time.Now()
		err := c.createNozzle(appGUID, instanceID)
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

func (c *Client) createNozzle(appGUID string, instanceID int) error {

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
