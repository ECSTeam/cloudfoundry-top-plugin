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

const NozzleInstances = 4

type Client struct {
	authToken     string
	options       *ClientOptions
	ui            terminal.UI
	cliConnection plugin.CliConnection
	eventrouting  *eventrouting.EventRouter
	router        *eventrouting.EventRouter
}

type ClientOptions struct {
	AppGUID        string
	Debug          bool
	Cygwin         bool
	NoFilter       bool
	Filter         string
	SubscriptionID string
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
	c.createNozzles(subscriptionID)

	ui.Start()

}

func (c *Client) createNozzles(subscriptionID string) {
	for i := 0; i < NozzleInstances; i++ {
		go c.createNozzle(subscriptionID, i)
	}
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

	c.routeEvents(messages, errors, instanceId)
	return nil
}

func (c *Client) routeEvents(messages <-chan *events.Envelope, errors <-chan error, instanceId int) {
	for {
		select {
		case envelope := <-messages:
			//debug.Debug(fmt.Sprintf("id: %v event:%v", instanceId, envelope))
			c.router.Route(envelope)
		case err := <-errors:
			c.handleError(err)
		}
	}
}

func (c *Client) handleError(err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		msg := fmt.Sprintf("Normal Websocket Closure: %v", err)
		debug.Error(msg)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		msg := fmt.Sprintf("Disconnected because nozzle couldn't keep up (CloseError): %v", err)
		debug.Error(msg)
	default:
		msg := fmt.Sprintf("Error reading firehose: %v", err)
		debug.Error(msg)
	}
}
