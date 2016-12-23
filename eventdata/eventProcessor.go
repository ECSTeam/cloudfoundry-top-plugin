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
	"fmt"
	"regexp"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
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

	ep := &EventProcessor{
		mu:               mu,
		cliConnection:    cliConnection,
		metadataManager:  metadataManager,
		startTime:        time.Now(),
		eventRateCounter: util.NewRateCounter(time.Second),
	}

	ep.currentEventData = NewEventData(mu, ep)
	ep.displayedEventData = NewEventData(mu, ep)

	return ep

}

func (ep *EventProcessor) Process(instanceId int, msg *events.Envelope) {
	ep.currentEventData.Process(instanceId, msg)
}

func (ep *EventProcessor) GetCliConnection() plugin.CliConnection {
	return ep.cliConnection
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
	processorCopy := ep.currentEventData.Clone()
	ep.displayedEventData = processorCopy
}

func (ep *EventProcessor) Start() {
	go ep.LoadCacheAndSeedData()
}

func (ep *EventProcessor) LoadCacheAndSeedData() {
	ep.metadataManager.LoadMetadata()
	ep.SeedStatsFromMetadata()
}

func (ep *EventProcessor) SeedStatsFromMetadata() {

	toplog.Info("EventProcessor>seedStatsFromMetadata")

	ep.mu.Lock()
	defer ep.mu.Unlock()

	ep.seedAppMap()
	ep.seedDomainMap()
	ep.seedRouteData()

	ep.currentEventData.EnableRouteTracking = true
}

func (ep *EventProcessor) seedAppMap() {

	currentStatsMap := ep.currentEventData.AppMap
	for _, app := range ep.metadataManager.GetAppMdManager().AllApps() {
		appId := app.Guid
		appStats := currentStatsMap[appId]
		if appStats == nil {
			// New app we haven't seen yet
			appStats = NewAppStats(appId)
			currentStatsMap[appId] = appStats
		}
	}
}

func (ep *EventProcessor) seedDomainMap() {
	currentStatsMap := ep.currentEventData.DomainMap
	for _, domain := range domain.AllDomains() {
		domainStats := currentStatsMap[domain.Name]
		if domainStats == nil {
			// New domain we haven't seen yet
			domainStats = NewDomainStats(domain.Guid)
			currentStatsMap[domain.Name] = domainStats
		}
	}
}

func (ep *EventProcessor) seedRouteData() {

	//currentDomainStatsMap := ep.currentEventData.DomainMap
	for _, route := range route.AllRoutes() {
		domainMd := domain.FindDomainMetadata(route.DomainGuid)
		ep.addRoute(domainMd.Name, route.Host, route.Path, route.Guid)
	}

	// Seed special host names
	apiDomain, apiHost := ep.getAPIHostAndDomain()
	ep.addRoute(apiDomain, apiHost, "", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/internal", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/internal/bulk/apps", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/internal/log_access", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/v2", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/v2/apps", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, apiHost, "/v2/syslog_drain_urls", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, "uaa", "", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, "uaa", "/oauth/token", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, "doppler", "", ep.generateUniqueRouteGuid())
	ep.addRoute(apiDomain, "doppler", "/apps", ep.generateUniqueRouteGuid())
	ep.addRoute("", "127.0.0.1", "/", ep.generateUniqueRouteGuid())
	ep.addRoute("", "127.0.0.1", "/v2", ep.generateUniqueRouteGuid())
	ep.addRoute("", "127.0.0.1", "/v2/stats/self", ep.generateUniqueRouteGuid())

	ep.addRoute(apiDomain, "proxy-0-p-mysql-ert", "/v0/backends", ep.generateUniqueRouteGuid())

}

func (ep *EventProcessor) generateUniqueRouteGuid() string {
	return util.Pseudo_uuid()
}

func (ep *EventProcessor) addRoute(domain, host, path, routeGuid string) {
	currentDomainStatsMap := ep.currentEventData.DomainMap

	//domainMd := metadata.FindDomainMetadataByName(domain)

	domainStats := currentDomainStatsMap[domain]
	if domainStats == nil {
		// New domain we haven't seen yet
		domainStats = NewDomainStats(domain)
		currentDomainStatsMap[domain] = domainStats
	}
	hostStats := domainStats.HostStatsMap[host]
	if hostStats == nil {
		hostStats = NewHostStats(host)
		domainStats.HostStatsMap[host] = hostStats
		//toplog.Info(fmt.Sprintf("seed hostStats: %v", host))
	}
	hostStats.AddPath(path, routeGuid)
}

func (ep *EventProcessor) getAPIHostAndDomain() (domain, host string) {
	apiUrl := util.GetApiEndpointNoProtocol(ep.cliConnection)

	parseInfoHostAndDomainNameStr := `^([^\.]+)\.([^\/^:]*)(?::[0-9]+)?`
	parseInfoHostAndDomainName := regexp.MustCompile(parseInfoHostAndDomainNameStr)
	parsedData := parseInfoHostAndDomainName.FindAllStringSubmatch(apiUrl, -1)
	if len(parsedData) != 1 {
		toplog.Debug(fmt.Sprintf("getAPIHostAndDomain>>Unable to parse (parsedData size) apiUri: %v", apiUrl))
		return
	}
	dataArray := parsedData[0]
	if len(dataArray) != 3 {
		toplog.Debug(fmt.Sprintf("getAPIHostAndDomain>>Unable to parse (dataArray size) apiUri: %v", apiUrl))
		return
	}
	host = dataArray[1]
	domain = dataArray[2]
	return domain, host
}

func (ep *EventProcessor) ClearStats() error {
	toplog.Info("EventProcessor>ClearStats")
	ep.currentEventData.Clear()
	ep.UpdateData()
	ep.SeedStatsFromMetadata()
	return nil
}
