package eventdata

import (
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/kkellner/cloudfoundry-top-plugin/metadata"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type EventProcessor struct {
	eventCount         int64
	eventRateCounter   *util.RateCounter
	eventRatePeak      int
	startTime          time.Time
	mu                 sync.Mutex
	currentEventData   *EventData
	displayedEventData *EventData
	cliConnection      plugin.CliConnection
}

func NewEventProcessor(cliConnection plugin.CliConnection) *EventProcessor {

	currentEventData := NewEventData()
	displayedEventData := NewEventData()

	return &EventProcessor{
		currentEventData:   currentEventData,
		displayedEventData: displayedEventData,
		cliConnection:      cliConnection,
		startTime:          time.Now(),
		eventRateCounter:   util.NewRateCounter(time.Second),
	}
}

func (ep *EventProcessor) Process(instanceId int, msg *events.Envelope) {
	ep.currentEventData.Process(instanceId, msg)
}

func (ep *EventProcessor) GetCurrentEventData() *EventData {
	return ep.currentEventData
}

func (ep *EventProcessor) GetDisplayedEventData() *EventData {
	return ep.displayedEventData
}

func (ep *EventProcessor) UpdateData() {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	processorCopy := ep.currentEventData.Clone()
	ep.displayedEventData = processorCopy
}

func (ep *EventProcessor) LoadMetadata() {
	toplog.Info("EventProcessor>loadMetadata")
	metadata.LoadAppCache(ep.cliConnection)
	metadata.LoadSpaceCache(ep.cliConnection)
	metadata.LoadOrgCache(ep.cliConnection)
}

func (ep *EventProcessor) Start() {
	go ep.LoadCacheAndSeeData()
}

func (ep *EventProcessor) LoadCacheAndSeeData() {
	ep.LoadMetadata()
	ep.SeedStatsFromMetadata()
}

func (ep *EventProcessor) SeedStatsFromMetadata() {

	toplog.Info("EventProcessor>seedStatsFromMetadata")

	ep.mu.Lock()
	defer ep.mu.Unlock()

	currentStatsMap := ep.currentEventData.AppMap
	for _, app := range metadata.AllApps() {
		appId := app.Guid
		appStats := currentStatsMap[appId]
		if appStats == nil {
			// New app we haven't seen yet
			appStats = NewAppStats(appId)
			currentStatsMap[appId] = appStats
		}
	}
}
