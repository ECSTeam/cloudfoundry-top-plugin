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
	currentProcessor   *AppStatsEventProcessor
	displayedProcessor *AppStatsEventProcessor
	cliConnection      plugin.CliConnection
}

func NewEventProcessor(cliConnection plugin.CliConnection) *EventProcessor {

	currentProcessor := NewAppStatsEventProcessor()
	displayedProcessor := NewAppStatsEventProcessor()

	return &EventProcessor{
		currentProcessor:   currentProcessor,
		displayedProcessor: displayedProcessor,
		cliConnection:      cliConnection,
		startTime:          time.Now(),
		eventRateCounter:   util.NewRateCounter(time.Second),
	}
}

func (self *EventProcessor) Process(instanceId int, msg *events.Envelope) {
	self.currentProcessor.Process(instanceId, msg)
}

func (self *EventProcessor) GetCurrentProcessor() *AppStatsEventProcessor {
	return self.currentProcessor
}

func (self *EventProcessor) GetDisplayedProcessor() *AppStatsEventProcessor {
	return self.displayedProcessor
}

func (self *EventProcessor) UpdateData() {
	self.mu.Lock()
	defer self.mu.Unlock()
	processorCopy := self.currentProcessor.Clone()
	self.displayedProcessor = processorCopy
}

func (self *EventProcessor) LoadMetadata() {
	toplog.Info("EventProcessor>loadMetadata")
	metadata.LoadAppCache(self.cliConnection)
	metadata.LoadSpaceCache(self.cliConnection)
	metadata.LoadOrgCache(self.cliConnection)
}

func (self *EventProcessor) Start() {
	go self.LoadCacheAndSeeData()
}

func (self *EventProcessor) LoadCacheAndSeeData() {
	self.LoadMetadata()
	self.SeedStatsFromMetadata()
}

func (asUI *EventProcessor) SeedStatsFromMetadata() {

	toplog.Info("AppStatsEventProcessor>seedStatsFromMetadata")

	asUI.mu.Lock()
	defer asUI.mu.Unlock()

	currentStatsMap := asUI.currentProcessor.AppMap
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
