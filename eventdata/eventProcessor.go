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
	"regexp"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"

	"strings"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type EventProcessor struct {
	eventCount         int64
	mu                 *sync.Mutex
	currentEventData   *EventData
	displayedEventData *EventData
	cliConnection      plugin.CliConnection
	privileged         bool
	metadataManager    *metadata.GlobalManager

	eventRateHistory *EventRateHistory

	eventRateCounterMap     map[events.Envelope_EventType]*util.RateCounter
	eventRateCounterMapLock sync.Mutex
}

func NewEventProcessor(cliConnection plugin.CliConnection, privileged bool) *EventProcessor {

	mu := &sync.Mutex{}

	metadataManager := metadata.NewGlobalManager(cliConnection)

	ep := &EventProcessor{
		mu:                  mu,
		cliConnection:       cliConnection,
		privileged:          privileged,
		metadataManager:     metadataManager,
		eventRateCounterMap: make(map[events.Envelope_EventType]*util.RateCounter),
	}

	ep.currentEventData = NewEventData(mu, ep)
	ep.displayedEventData = NewEventData(mu, ep)
	ep.eventRateHistory = NewEventRateHistory(ep)
	ep.eventRateHistory.start()
	return ep

}

func (ep *EventProcessor) Process(instanceId int, msg *events.Envelope) {

	eventType := msg.GetEventType()
	ep.eventRateCounterMapLock.Lock()
	eventCounter := ep.eventRateCounterMap[eventType]
	if eventCounter == nil {
		eventCounter = util.NewRateCounter(time.Second)
		ep.eventRateCounterMap[eventType] = eventCounter
	}
	ep.eventRateCounterMapLock.Unlock()
	eventCounter.Incr()
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

func (ep *EventProcessor) GetCurrentEventRateHistory() *EventRateHistory {
	return ep.eventRateHistory
}

func (ep *EventProcessor) GetMetadataManager() *metadata.GlobalManager {
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
	if ep.privileged {
		ep.seedSpecialRouteData()
	}

	ep.currentEventData.EnableRouteTracking = true
}

func (ep *EventProcessor) seedAppMap() {

	currentStatsMap := ep.currentEventData.AppMap
	for _, app := range ep.metadataManager.GetAppMdManager().AllApps() {
		appId := app.Guid
		appStats := currentStatsMap[appId]
		if appStats == nil {
			// New app we haven't seen yet
			appStats = eventApp.NewAppStats(appId)
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
			domainStats = eventRoute.NewDomainStats(domain.Guid)
			currentStatsMap[strings.ToLower(domain.Name)] = domainStats
		}
	}
}

func (ep *EventProcessor) seedRouteData() {

	for _, route := range route.AllRoutes() {
		domainMd := domain.FindDomainMetadata(route.DomainGuid)
		ep.addRoute(domainMd.Name, route.Host, route.Path, route.Port, route.Guid)
	}
}

func (ep *EventProcessor) seedSpecialRouteData() {

	// Seed special host names
	apiDomain, apiHost := ep.getAPIHostAndDomain()

	// Register all documented APIs from https://apidocs.cloudfoundry.org/249/
	// TODO: Add a way to dynamically add level-n paths as seen during runtime
	// I.e., pathLevel=2 on "api" host would track: /*/*
	// E.g., /v2/apps or /v3/jobs/2323-232-2323 (only /v2/jobs would be tracked)
	registerApiPaths := [...]string{
		"/internal",
		"/internal/bulk/apps",
		"/internal/log_access",
		"/v2",
		"/v2/app_usage_events",
		"/v2/apps",
		"/v2/buildpacks",
		"/v2/domains",
		"/v2/config",
		"/v2/events",
		"/v2/info",
		"/v2/jobs",
		"/v2/quota_definitions",
		"/v2/organizations",
		"/v2/private_domains",
		"/v2/resource_match",
		"/v2/routes",
		"/v2/route_mappings",
		"/v2/security_groups",
		"/v2/service_bindings",
		"/v2/service_brokers",
		"/v2/service_instances",
		"/v2/service_keys",
		"/v2/service_plan_visibilities",
		"/v2/service_plans",
		"/v2/service_usage_events",
		"/v2/services",
		"/v2/shared_domains",
		"/v2/space_quota_definitions",
		"/v2/spaces",
		"/v2/stacks",
		"/v2/user_provided_service_instances",
		"/v2/users",
		"/v2/syslog_drain_urls",
	}

	for _, registerPath := range registerApiPaths {
		ep.addInternalRoute(apiDomain, apiHost, registerPath, 0)
	}

	ep.addInternalRoute(apiDomain, "uaa", "", 0)
	ep.addInternalRoute(apiDomain, "uaa", "/oauth/token", 0)
	ep.addInternalRoute(apiDomain, "doppler", "", 0)
	ep.addInternalRoute(apiDomain, "doppler", "/apps", 0)

}

func (ep *EventProcessor) addInternalRoute(domainName string, hostName string, pathName string, port int) *eventRoute.RouteStats {
	domain := domain.FindDomainMetadataByName(domainName)
	if domain == nil {
		toplog.Warn("addInternalRoute for domain [%v] but domain metadata not found (or not loaded) host: [%v] path: [%v]", domainName, hostName, pathName)
		return nil
	}
	route := route.CreateInternalGeneratedRoute(hostName, pathName, domain.Guid, port)
	return ep.addRoute(domainName, hostName, pathName, port, route.Guid)
}

func (ep *EventProcessor) addRoute(domain string, host string, path string, port int, routeGuid string) *eventRoute.RouteStats {

	domain = strings.ToLower(domain)
	host = strings.ToLower(host)

	currentDomainStatsMap := ep.currentEventData.DomainMap
	domainStats := currentDomainStatsMap[domain]
	if domainStats == nil {
		// New domain we haven't seen yet
		domainStats = eventRoute.NewDomainStats(domain)
		currentDomainStatsMap[domain] = domainStats
	}
	hostStats := domainStats.HostStatsMap[host]
	if hostStats == nil {
		hostStats = eventRoute.NewHostStats(host)
		domainStats.HostStatsMap[host] = hostStats
		//toplog.Info("seed hostStats: %v", host)
	}
	if port == 0 {
		return hostStats.AddPath(path, routeGuid, true)
	} else {
		return hostStats.AddPort(port, routeGuid, true)
	}
}

func (ep *EventProcessor) getAPIHostAndDomain() (domain, host string) {
	apiUrl := util.GetApiEndpointNoProtocol(ep.cliConnection)

	parseInfoHostAndDomainNameStr := `^([^\.]+)\.([^\/^:]*)(?::[0-9]+)?`
	parseInfoHostAndDomainName := regexp.MustCompile(parseInfoHostAndDomainNameStr)
	parsedData := parseInfoHostAndDomainName.FindAllStringSubmatch(apiUrl, -1)
	if len(parsedData) != 1 {
		toplog.Debug("getAPIHostAndDomain>>Unable to parse (parsedData size) apiUri: %v", apiUrl)
		return
	}
	dataArray := parsedData[0]
	if len(dataArray) != 3 {
		toplog.Debug("getAPIHostAndDomain>>Unable to parse (dataArray size) apiUri: %v", apiUrl)
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
