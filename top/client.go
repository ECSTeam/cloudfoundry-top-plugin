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

	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/ui"
	gops "github.com/mitchellh/go-ps"
)

type Client struct {
	authToken     string
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
		return
	}

	//isDebug := c.options.Debug
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

	c.authToken, err = conn.AccessToken()
	if err != nil {
		c.ui.Failed(err.Error())
		return
	}

	ui := ui.NewMasterUI(conn)
	c.router = ui.GetRouter()

	toplog.Info("Top started at " + time.Now().Format("01-02-2006 15:04:05"))

	subscriptionID := "TopPlugin_" + pseudo_uuid()
	go c.createNozzles(subscriptionID)

	ui.Start()

}

func (c *Client) createNozzles(subscriptionID string) {
	for i := 0; i < c.options.Nozzles; i++ {
		go c.createAndKeepAliveNozzle(subscriptionID, i)
	}
	toplog.Info(fmt.Sprintf("Starting %v nozzle instances", c.options.Nozzles))
}

func (c *Client) createAndKeepAliveNozzle(subscriptionID string, instanceId int) error {

	minRetrySeconds := (2 * time.Second)

	for {
		// This is a blocking call if no error
		startTime := time.Now()
		err := c.createNozzle(subscriptionID, instanceId)
		if err != nil {
			toplog.Error(fmt.Sprintf("Nozzle #%v - Stopped with error: %v", instanceId, err))
			break
		}
		toplog.Warn(fmt.Sprintf("Nozzle #%v - Shutdown. Nozzle instance will be restarted", instanceId))
		lastRetry := time.Now().Sub(startTime)
		if lastRetry < minRetrySeconds {
			toplog.Info(fmt.Sprintf("Nozzle #%v - Nozzle instance restart too fast, delaying for %v", instanceId, minRetrySeconds))
			time.Sleep(minRetrySeconds)
		}
	}
	return nil
}

func (c *Client) createNozzle(subscriptionID string, instanceId int) error {

	conn := c.cliConnection

	dopplerEndpoint, err := conn.DopplerEndpoint()
	if err != nil {
		//c.ui.Failed(err.Error())
		return err
	}

	skipVerifySSL, err := conn.IsSSLDisabled()
	if err != nil {
		//c.ui.Failed("couldn't check if ssl verification is disabled: " + err.Error())
		return err
	}

	dopplerConnection := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipVerifySSL}, nil)

	dopplerConnection.SetMinRetryDelay(500 * time.Millisecond)
	dopplerConnection.SetMaxRetryDelay(15 * time.Second)
	dopplerConnection.SetIdleTimeout(15 * time.Second)

	messages, errors := dopplerConnection.Firehose(subscriptionID, c.authToken)
	defer dopplerConnection.Close()

	toplog.Info(fmt.Sprintf("Nozzle #%v - Started", instanceId))

	c.routeEvents(instanceId, messages, errors)
	return nil
}

func (c *Client) routeEvents(instanceId int, messages <-chan *events.Envelope, errors <-chan error) {
	for {
		select {
		case envelope := <-messages:
			c.router.Route(instanceId, envelope)
		case err := <-errors:
			c.handleError(instanceId, err)
			// Nozzle connection does not seem to recover from errors well, so
			// return here so it can be closed and a new instanced opened
			return
		}
	}
}

func (c *Client) handleError(instanceId int, err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		msg := fmt.Sprintf("Nozzle #%v - Normal Websocket Closure: %v", instanceId, err)
		toplog.Error(msg)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		msg := fmt.Sprintf("Nozzle #%v - Disconnected because nozzle couldn't keep up (CloseError): %v", instanceId, err)
		toplog.Error(msg)
	default:
		msg := fmt.Sprintf("Nozzle #%v - Error reading firehose: %v", instanceId, err)
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
