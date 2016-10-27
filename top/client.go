package top


import (

	"fmt"
	"crypto/tls"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/gorilla/websocket"

	"github.com/kkellner/cloudfoundry-top-plugin/eventrouting"
	"github.com/kkellner/cloudfoundry-top-plugin/appStats"

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
	NoFilter       bool
	Filter         string
	SubscriptionID string
}


var (
	dopplerConnection *consumer.Consumer
	apiEndpoint string
	appStatsUI *appStats.AppStatsUI
)

func NewClient(cliConnection plugin.CliConnection, options *ClientOptions, ui terminal.UI) *Client {

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		ui.Failed(err.Error())
	}


	authToken, err := cliConnection.AccessToken()
	if err != nil {
		ui.Failed(err.Error())
	}

	apiEndpoint, err = cliConnection.ApiEndpoint()
	if err != nil {
		ui.Failed(err.Error())
	}

	return &Client{
		dopplerEndpoint: dopplerEndpoint,
		authToken:       authToken,
		options:         options,
		ui:              ui,
		cliConnection:   cliConnection,
	}

}



func (c *Client) Start() {

	dopplerConnection = consumer.New(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	if c.options.Debug {
		dopplerConnection.SetDebugPrinter(ConsoleDebugPrinter{ui: c.ui})
	}

	subscriptionID := "TopPlugin_" + pseudo_uuid()
	c.ui.Say("Starting the nozzle for monitoring.  subscriptionID:"+subscriptionID)

	// consumer.Stream(appGuid, authToken)
	// consumer.ContainerEnvelopes(appId, authToken)
	// consumer.Firehose(appId, authToken)
	//c.messages, c.errors = dopplerConnection.FirehoseWithoutReconnect(subscriptionID, c.authToken)
	c.messages, c.errors = dopplerConnection.Firehose(subscriptionID, c.authToken)

	defer dopplerConnection.Close()
	c.ui.Say("Hit Ctrl+c to exit")


	//c.router = eventrouting.NewEventRouter(appStatsUI.GetProcessor())
	go c.routeEvent()

	ui := NewMasterUI(c.cliConnection)
	c.router = ui.GetRouter()

	ui.Start()
	dopplerConnection.Close()
}


func (c *Client) routeEvent() error {

	for {
		select {
		case envelope := <-c.messages:
			c.router.Route(envelope)
		case err := <-c.errors:
			c.handleError(err)
			return err
		}
	}
}

func (c *Client) handleError(err error) {

	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		fmt.Printf("Normal Websocket Closure: %v", err)
	case websocket.IsCloseError(err, websocket.ClosePolicyViolation):
		fmt.Printf("Error while reading from the firehose: %v", err)
		fmt.Printf("Disconnected because nozzle couldn't keep up. Please try scaling up the nozzle.", nil)
	default:
		fmt.Printf("Error while reading from the firehose: %v", err)
	}
	fmt.Printf("Closing connection with traffic controller due to %v", err)
	dopplerConnection.Close()

}


type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}
