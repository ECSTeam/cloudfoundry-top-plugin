package top


import (

	"fmt"
	"crypto/tls"
	"time"
	//"errors"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/gorilla/websocket"

	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"

)

type Client struct {
	dopplerEndpoint string
	authToken       string
	options         *ClientOptions
	ui              terminal.UI
	cliConnection   plugin.CliConnection
	eventrouting 		*eventrouting.EventRouter
	errors 					<-chan error
	messages 				<-chan *events.Envelope
	router					*eventrouting.EventRouter
}

type ClientOptions struct {
	AppGUID        string
	Debug          bool
	Cygwin				 bool
	NoFilter       bool
	Filter         string
	SubscriptionID string
}


var (
	dopplerConnection *consumer.Consumer
	apiEndpoint string
	//appStatsUI *appStats.AppStatsUI
)

func NewClient(cliConnection plugin.CliConnection, options *ClientOptions, ui terminal.UI) *Client {

	return &Client{
		options:         options,
		ui:              ui,
		cliConnection:   cliConnection,
	}

}

func (c *Client) Start() {

  isDebug := c.options.Debug
	conn := c.cliConnection

	isLoggedIn, err := conn.IsLoggedIn()
	if err !=nil {
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

	c.dopplerEndpoint, err = conn.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
		return
	}

	apiEndpoint, err = conn.ApiEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
		return
	}

	skipVerifySSL, err := conn.IsSSLDisabled()
	if err != nil {
		c.ui.Failed("couldn't check if ssl verification is disabled: " + err.Error())
		return
	}

	subscriptionID := "TopPlugin_" + pseudo_uuid()
	dopplerConnection = consumer.New(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: skipVerifySSL}, nil)
	if isDebug {
		//dopplerConnection.SetDebugPrinter(ConsoleDebugPrinter{ui: c.ui})
		c.ui.Say("Starting the nozzle for monitoring.  subscriptionID:"+subscriptionID)
		c.ui.Say("Hit Ctrl+c to exit")
	}

	dopplerConnection.SetMinRetryDelay(500 * time.Millisecond)
	dopplerConnection.SetMaxRetryDelay(15 * time.Second)
	dopplerConnection.SetIdleTimeout(15 * time.Second)

	c.messages, c.errors = dopplerConnection.Firehose(subscriptionID, c.authToken)
	defer dopplerConnection.Close()

	ui := NewMasterUI(c.cliConnection)
	c.router = ui.GetRouter()

	debug.Info("Top started at "+time.Now().Format("01-02-2006 15:04:05"))

	go c.routeEvent()
	ui.Start()

}


func (c *Client) routeEvent() error {

	for {
		select {
		case envelope := <-c.messages:
			c.router.Route(envelope)
		case err := <-c.errors:
			c.handleError(err)
			//return err
		}
	}
}

func (c *Client) handleError(err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		msg := fmt.Sprintf("Normal Websocket Closure: %v", err)
		//fmt.Printf(msg)
		debug.Error(msg)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		msg := fmt.Sprintf("Error while reading from the firehose: %v", err)
		debug.Error(msg)
		//fmt.Printf("Disconnected because nozzle couldn't keep up. Please try scaling up the nozzle.", nil)
	default:
		msg := fmt.Sprintf("Error while reading from the firehose: %v", err)
		debug.Error(msg)
	}
	//fmt.Printf("Closing connection with traffic controller due to %v", err)
	//dopplerConnection.Close()

}


type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
