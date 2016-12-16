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

package eventdata

import (
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type EventProcessor struct {
	eventCount         int64
	eventRateCounter   *util.RateCounter
	eventRatePeak      int
	startTime          time.Time
	mu                 *sync.Mutex
	currentEventData   *EventData
	displayedEventData *EventData
	cliConnection      plugin.CliConnection
	metadataManager    *metadata.Manager
}

func NewEventProcessor(cliConnection plugin.CliConnection) *EventProcessor {

	mu := &sync.Mutex{}

	metadataManager := metadata.NewManager(cliConnection)

	currentEventData := NewEventData(mu, metadataManager)
	displayedEventData := NewEventData(mu, metadataManager)

	ep := &EventProcessor{
		mu:                 mu,
		currentEventData:   currentEventData,
		displayedEventData: displayedEventData,
		cliConnection:      cliConnection,
		metadataManager:    metadataManager,
		startTime:          time.Now(),
		eventRateCounter:   util.NewRateCounter(time.Second),
	}
	return ep

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

func (ep *EventProcessor) GetMetadataManager() *metadata.Manager {
	return ep.metadataManager
}

func (ep *EventProcessor) UpdateData() {
	//ep.mu.Lock()
	//defer ep.mu.Unlock()
	processorCopy := ep.currentEventData.Clone()
	ep.displayedEventData = processorCopy
}

func (ep *EventProcessor) Start() {
	go ep.LoadCacheAndSeeData()
}

func (ep *EventProcessor) LoadCacheAndSeeData() {
	ep.metadataManager.LoadMetadata()
	ep.SeedStatsFromMetadata()
}

func (ep *EventProcessor) SeedStatsFromMetadata() {

	toplog.Info("EventProcessor>seedStatsFromMetadata")

	ep.mu.Lock()
	defer ep.mu.Unlock()

	currentStatsMap := ep.currentEventData.AppMap
	for _, app := range ep.metadataManager.GetAppMdManager().AllApps() {
		appId := app.Guid
		appStats := currentStatsMap[appId]
		if appStats == nil {
			// New app we haven't seen yet
			appStats = NewAppStats(appId)
			currentStatsMap[appId] = appStats // Thread was here at crash
		}
	}
}

func (ep *EventProcessor) ClearStats() error {
	toplog.Info("EventProcessor>ClearStats")
	ep.currentEventData.Clear()
	ep.UpdateData()
	ep.SeedStatsFromMetadata()
	return nil
}
