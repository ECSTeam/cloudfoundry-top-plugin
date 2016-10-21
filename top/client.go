package top


import (
	"crypto/tls"
	"strconv"

	"fmt"
	"time"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/cloudfoundry/sonde-go/events"

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
	ctr = 0

	//appMap map[UUIDKey]int
	appMap = make(map[UUIDKey]int)

	dopplerConnection *consumer.Consumer
)

func NewClient(authToken, doppplerEndpoint string, options *ClientOptions, ui terminal.UI) *Client {
	return &Client{
		dopplerEndpoint: doppplerEndpoint,
		authToken:       authToken,
		options:         options,
		ui:              ui,
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
  lookupUUIDKey := &UUIDKey{0, 0}

	//go say(appMap, "world")

	for envelope := range output {

		// Check if this is an HttpStartStop event
		if filter == "" || filter == strconv.Itoa((int)(envelope.GetEventType())) {
			appId := envelope.GetHttpStartStop().GetApplicationId()
			instId := envelope.GetHttpStartStop().GetInstanceId()

			// Check if this is an application event
			if appId != nil && instId != "" {
				lookupUUIDKey.Low = appId.GetLow()
				lookupUUIDKey.High = appId.GetHigh()
				count := appMap[*lookupUUIDKey]
				count++
				appMap[*lookupUUIDKey] = count
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

	if v, err := g.SetView("helloView", maxX/2-32, maxY/2, maxX/2+32, maxY/2+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//fmt.Fprintln(v, "Hello world!")
		fmt.Fprintln(v, "waiting...")
	}

	if v, err := g.SetView("summaryView", 0, 0, maxX-1, 8); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Summary"
			v.Frame = true
			fmt.Fprintln(v, "View #2")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	dopplerConnection.Close()
	return gocui.ErrQuit
}


func counter(g *gocui.Gui) {
	//defer wg.Done()

	for {
		select {
		case <-doneX:
			return
		case <-time.After(1000 * time.Millisecond):
			mu.Lock()
			m := appMap
			ctr++
			mu.Unlock()

			g.Execute(func(g *gocui.Gui) error {
				v, err := g.View("helloView")
				if err != nil {
					return err
				}
				if len(appMap) > 0 {
					v.Clear()
					//fmt.Fprintln(v, n)
					for appId, count := range m {
						//fmt.Println(s)
						fmt.Fprintf(v, "%v count:%d\n", appId, count)
					}
				}

				v, err = g.View("summaryView")
				if err != nil {
					return err
				}
				v.Clear()
				fmt.Fprintf(v, "Unique Apps:%5v  ", len(appMap))
				fmt.Fprintf(v, "More stats:%v", len(appMap))
				return nil
			})
		}
	}
}
