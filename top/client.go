package top


import (
	"crypto/tls"
	"strconv"

	"fmt"
	"time"
	"encoding/binary"
	"strings"
	"encoding/json"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry/cli/plugin"

  cfclient "github.com/cloudfoundry-community/go-cfclient"

	"github.com/jroimartin/gocui"
	"log"
	"sync"
	//"syscall"

)

type Client struct {
	dopplerEndpoint string
	authToken       string
	options         *ClientOptions
	ui              terminal.UI
	cliConnection   plugin.CliConnection
}

type ClientOptions struct {
	AppGUID        string
	Debug          bool
	NoFilter       bool
	Filter         string
	SubscriptionID string
}

type UUIDKey struct {
	Low              uint64
	High             uint64
}

var (
	doneX = make(chan bool)
	wg   sync.WaitGroup

	mu  sync.Mutex // protects ctr

	//appMap = make(map[UUIDKey]int)
	appMap = make(map[string]int)
	totalEvents = 0

	dopplerConnection *consumer.Consumer
	appsMetadata []cfclient.App
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

	return &Client{
		dopplerEndpoint: dopplerEndpoint,
		authToken:       authToken,
		options:         options,
		ui:              ui,
		cliConnection:   cliConnection,
	}

}

func GetAppMetadata(c *Client) {


		//requestUrl := "/v2/apps?inline-relations-depth=2"
		requestUrl := "/v2/apps"
		reponseJSON, err := c.cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
		if err != nil {
			c.ui.Failed(err.Error())
			return
		}

		var appResp cfclient.AppResponse
		// joining since it's an array of strings
		outputStr := strings.Join(reponseJSON, "")
		fmt.Printf("Response Size: %v\n", len(outputStr))
		outputBytes := []byte(outputStr)
		err2 := json.Unmarshal(outputBytes, &appResp)
		if err2 != nil {
					c.ui.Failed(err.Error())
		}

		//var apps []cfclient.App
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
			app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
			appsMetadata = append(appsMetadata, app.Entity)
		}

		for _, app := range appsMetadata {
			fmt.Printf("appName: %v  %v\n", app.Name, app.Guid)
		}


}

func (c *Client) Start() {

	var err error
	dopplerConnection = consumer.New(c.dopplerEndpoint, &tls.Config{InsecureSkipVerify: true}, nil)
	if c.options.Debug {
		dopplerConnection.SetDebugPrinter(ConsoleDebugPrinter{ui: c.ui})
	}
	filter := ""
	switch {
	case c.options.NoFilter:
		filter = ""
	case c.options.Filter != "":
		envelopeType, ok := events.Envelope_EventType_value[c.options.Filter]
		if !ok {
			c.ui.Warn("Unable to recognize filter %s", c.options.Filter)
			return
		}
		filter = strconv.Itoa(int(envelopeType))

	default:
		//c.ui.Say("What type of firehose messages do you want to see?")
		filter, err = c.promptFilterType()
		if err != nil {
			c.ui.Warn(err.Error())
			return
		}
	}


	GetAppMetadata(c)



	var errors <-chan error
	var output <-chan *events.Envelope
	if len(c.options.AppGUID) != 0 {
		c.ui.Say("Starting the nozzle for app %s", c.options.AppGUID)
		output, errors = dopplerConnection.StreamWithoutReconnect(c.options.AppGUID, c.authToken)
	} else {
		subscriptionID := c.options.SubscriptionID
		if len(subscriptionID) == 0 {
			subscriptionID = "TopPlugin"
		}
		c.ui.Say("Starting the nozzle for monitoring")
		output, errors = dopplerConnection.FirehoseWithoutReconnect(subscriptionID, c.authToken)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for err := range errors {
			c.ui.Warn(err.Error())
			return
		}
	}()


	defer dopplerConnection.Close()

	c.ui.Say("Hit Ctrl+c to exit")

  if output == nil || filter == "" {
		c.ui.Say("whatever")
	}

// *******************

	go initGui()

	//var appMap map[UUIDKey]int
	//appMap = make(map[UUIDKey]int)

	// Create once outside loop
  //lookupUUIDKey := &UUIDKey{0, 0}

	//go say(appMap, "world")

	for envelope := range output {
		totalEvents++

		// Took the following from the nozzle code -- how does it know that WE are the slow one and not another nozzle?
		if envelope.GetEventType() == events.Envelope_CounterEvent && envelope.CounterEvent.GetName() == "TruncatingBuffer.DroppedMessages" && envelope.GetOrigin() == "doppler" {
			c.ui.Say("We've intercepted an upstream message which indicates that the nozzle or the TrafficController is not keeping up. Please try scaling up the nozzle.")
		}

		//time.Sleep(1000 * time.Millisecond)

		// Check if this is an HttpStartStop event
		if filter == "" || filter == strconv.Itoa((int)(envelope.GetEventType())) {

			appUUID := envelope.GetHttpStartStop().GetApplicationId()
			instId := envelope.GetHttpStartStop().GetInstanceId()

			// Check if this is an application event
			if appUUID != nil && instId != "" {

				appId := formatUUID(appUUID)
				//c.ui.Say("**** appId:%v ****", appId)

				count := appMap[appId]
				count++
				appMap[appId] = count
				//c.ui.Say("%v size:%d  count:%d\n", appId, len(appMap), count)

				//if envelope.GetHttpStartStop().GetPeerType() == events.PeerType_Client {
				//	c.ui.Say("CLIENT EVENT \n")
				//}

				//c.ui.Say("%v \n", envelope)
			}
		}
	}

	/*
	c.ui.Say("after for envelope loop")
	for error := range errors {
		c.ui.Say("ERROR event from top: %v \n", error)
	}
	<-done

	*/
}

func say(appMap map[UUIDKey]int, s string) {
	//for i := 0; i < 50; i++ {
	for {
		for appId, count := range appMap {
			//fmt.Println(s)
			fmt.Printf("%v size:%d  count:%d\n", appId, len(appMap), count)
		}
		fmt.Printf("-\n")
		time.Sleep(1000 * time.Millisecond)
	}
}

func (c *Client) promptFilterType() (string, error) {

  filter := "4"
	/*
	filter := c.ui.Ask(`Please enter one of the following choices:
	  hit 'enter' for all messages
	  2 for HttpStart
	  3 for HttpStop
	  4 for HttpStartStop
	  5 for LogMessage
	  6 for ValueMetric
	  7 for CounterEvent
	  8 for Error
	  9 for ContainerMetric
	`)
  */
	if filter == "" {
		return "", nil
	}

	filterInt, err := strconv.Atoi(filter)
	if err != nil {
		return "", fmt.Errorf("Invalid filter choice %s. Enter an index from 2-9", filter)
	}

	_, ok := events.Envelope_EventType_name[int32(filterInt)]
	if !ok {
		return "", fmt.Errorf("Invalid filter choice %d", filterInt)
	}

	return filter, nil
}

type ConsoleDebugPrinter struct {
	ui terminal.UI
}

func (p ConsoleDebugPrinter) Print(title, dump string) {
	p.ui.Say(title)
	p.ui.Say(dump)
}

func initGui() {

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'Q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'c', gocui.ModNone, clearStats); err != nil {
		log.Panicln(err)
	}
	//wg.Add(1)
	go counter(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	// Send this process a SIGHUP
  //go syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
  //waitSig(t, c, syscall.SIGHUP)

}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	//if v, err := g.SetView("helloView", maxX/2-32, maxY/2, maxX/2+32, maxY/2+4); err != nil {
	if v, err := g.SetView("detailView", 0, 5, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//fmt.Fprintln(v, "Hello world!")
		fmt.Fprintln(v, "")
	}

	if v, err := g.SetView("summaryView", 0, 0, maxX-1, 4); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Summary"
			v.Frame = true
			fmt.Fprintln(v, "")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	dopplerConnection.Close()
	return gocui.ErrQuit
}

func clearStats(g *gocui.Gui, v *gocui.View) error {
	appMap = make(map[string]int)
	totalEvents = 0
	updateDisplay(g)
	return nil
}

func counter(g *gocui.Gui) {
	for {
		select {
		case <-doneX:
			return
		case <-time.After(1000 * time.Millisecond):
			updateDisplay(g)
		}
	}
}

func updateDisplay(g *gocui.Gui) {
	mu.Lock()
	m := appMap
	mu.Unlock()

	//maxX, maxY := g.Size()

	g.Execute(func(g *gocui.Gui) error {
		v, err := g.View("detailView")
		if err != nil {
			return err
		}
		if len(m) > 0 {
			v.Clear()
			fmt.Fprintf(v, "%-40v %6v %6v %6v %6v %6v\n", "Application","2xx","3xx","4xx","5xx","Total")
			for appId, count := range m {
				appName := findAppName(appId)
				fmt.Fprintf(v, "%-40v %6d %6d %6d %6d %6d\n", appName, 0,0,0 ,0,count)
			}
		} else {
			v.Clear()
			fmt.Fprintln(v, "No data yet...")
		}

		v, err = g.View("summaryView")
		if err != nil {
			return err
		}
		v.Clear()
		/*
		for i := 0; i < 200000; i++ {
			v.Clear()
		}
		*/
		fmt.Fprintf(v, "Total events: %-11v", totalEvents)
		fmt.Fprintf(v, "Unique Apps: %-11v", len(m))
		fmt.Fprintf(v, "%v\n", time.Now().Format("2006-01-02 15:04:05.000"))

		return nil
	})
}

func findAppName(appId string) string {
	for _, app := range appsMetadata {
		if app.Guid == appId {
			return app.Name;
		}
	}
	return appId
}

func formatUUID(uuid *events.UUID) string {
	var uuidBytes [16]byte
	binary.LittleEndian.PutUint64(uuidBytes[:8], uuid.GetLow())
	binary.LittleEndian.PutUint64(uuidBytes[8:], uuid.GetHigh())
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])
}
