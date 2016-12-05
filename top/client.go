package top

import (
	"crypto/tls"
	"fmt"
	"time"
	//"errors"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gorilla/websocket"

	"github.com/kkellner/cloudfoundry-top-plugin/debug"
	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
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
	AppGUID string
	Debug   bool
	Cygwin  bool
	Nozzles int
}

func NewClient(cliConnection plugin.CliConnection, options *ClientOptions, ui terminal.UI) *Client {

	return &Client{
		options:       options,
		ui:            ui,
		cliConnection: cliConnection,
	}
}

func (c *Client) Start() {

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

	ui := NewMasterUI(conn)
	c.router = ui.GetRouter()

	debug.Info("Top started at " + time.Now().Format("01-02-2006 15:04:05"))

	subscriptionID := "TopPlugin_" + pseudo_uuid()
	go c.createNozzles(subscriptionID)

	ui.Start()

}

func (c *Client) createNozzles(subscriptionID string) {
	for i := 0; i < c.options.Nozzles; i++ {
		go c.createNozzle(subscriptionID, i)
	}
	debug.Info(fmt.Sprintf("Created %v nozzles", c.options.Nozzles))
}

func (c *Client) createNozzle(subscriptionID string, instanceId int) error {

	conn := c.cliConnection

	dopplerEndpoint, err := conn.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
		return err
	}

	skipVerifySSL, err := conn.IsSSLDisabled()
	if err != nil {
		c.ui.Failed("couldn't check if ssl verification is disabled: " + err.Error())
		return err
	}

	dopplerConnection := consumer.New(dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipVerifySSL}, nil)

	dopplerConnection.SetMinRetryDelay(500 * time.Millisecond)
	dopplerConnection.SetMaxRetryDelay(15 * time.Second)
	dopplerConnection.SetIdleTimeout(15 * time.Second)

	messages, errors := dopplerConnection.Firehose(subscriptionID, c.authToken)
	defer dopplerConnection.Close()

	c.routeEvents(instanceId, messages, errors)
	return nil
}

func (c *Client) routeEvents(instanceId int, messages <-chan *events.Envelope, errors <-chan error) {
	for {
		select {
		case envelope := <-messages:
			//debug.Debug(fmt.Sprintf("id: %v event:%v", instanceId, envelope))
			c.router.Route(instanceId, envelope)
		case err := <-errors:
			c.handleError(instanceId, err)
		}
	}
}

func (c *Client) handleError(instanceId int, err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		msg := fmt.Sprintf("Nozzle #%v - Normal Websocket Closure: %v", instanceId, err)
		debug.Error(msg)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		msg := fmt.Sprintf("Nozzle #%v - Disconnected because nozzle couldn't keep up (CloseError): %v", instanceId, err)
		debug.Error(msg)
	default:
		msg := fmt.Sprintf("Nozzle #%v - Error reading firehose: %v", instanceId, err)
		debug.Error(msg)
	}
}
