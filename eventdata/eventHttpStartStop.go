// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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
	"strings"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/config"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventRoute"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

func (ed *EventData) httpStartStopEvent(msg *events.Envelope) {

	appUUID := msg.GetHttpStartStop().GetApplicationId()
	instId := msg.GetHttpStartStop().GetInstanceId()
	//instIndex := msg.GetHttpStartStop().GetInstanceIndex()
	httpEvent := msg.GetHttpStartStop()
	peerType := httpEvent.GetPeerType()
	switch {
	case peerType == events.PeerType_Client:

		if ed.EnableRouteTracking {
			ed.handleRouteStats(msg)
		}

		if appUUID != nil && instId != "" {
			ed.httpStartStopEventForApp(msg)
		} else if appUUID == nil && instId != "" {
			switch httpEvent.GetMethod() {
			case events.Method_PUT:
				fallthrough
			//case events.Method_POST:
			// TODO: Maybe this shouldn't be here as POST doesn't seem to include guid on URI
			//	fallthrough
			case events.Method_DELETE:
				// Check if we have a PCF API call
				uri := httpEvent.GetUri()

				isApiCall, apiUri := ed.checkIfApiCall_PCF1_6(uri)
				if !isApiCall {
					// PCF 1.7 needs more testing -- comment out this check for now
					//isApiCall, apiUri = ed.checkIfApiCall_PCF1_7(uri)
				}
				if isApiCall {
					ed.pcfApiHasBeenCalled(msg, apiUri)
				}
			}
		}
	default:
		// Ignore
	}

}

func (ed *EventData) httpStartStopEventForApp(msg *events.Envelope) {

	appUUID := msg.GetHttpStartStop().GetApplicationId()
	instId := msg.GetHttpStartStop().GetInstanceId()
	//instIndex := msg.GetHttpStartStop().GetInstanceIndex()
	httpStartStopEvent := msg.GetHttpStartStop()

	//toplog.Debug("index: %v\n", instIndex)
	//toplog.Debug("index mem: %v\n", msg.GetHttpStartStop().InstanceIndex)
	//fmt.Printf("index: %v\n", instIndex)
	ed.TotalEvents++
	appId := formatUUID(appUUID)
	//c.ui.Say("**** appId:%v ****", appId)

	appStats := ed.getAppStats(appId)
	if appStats.AppUUID == nil {
		appStats.AppUUID = appUUID
	}

	containerTraffic := ed.getContainerTraffic(appStats, instId)

	responseTimeNano := *httpStartStopEvent.StopTimestamp - *httpStartStopEvent.StartTimestamp

	containerTraffic.ResponseL60Time.Track(responseTimeNano)
	containerTraffic.ResponseL10Time.Track(responseTimeNano)
	containerTraffic.ResponseL1Time.Track(responseTimeNano)

	statusCode := httpStartStopEvent.GetStatusCode()
	httpMethod := httpStartStopEvent.GetMethod()

	httpStatusCodeMap := containerTraffic.HttpInfoMap[httpMethod]
	if httpStatusCodeMap == nil {
		httpStatusCodeMap = make(map[int32]*eventApp.HttpInfo)
		containerTraffic.HttpInfoMap[httpMethod] = httpStatusCodeMap
	}

	httpInfo := httpStatusCodeMap[statusCode]
	if httpInfo == nil {
		httpInfo = eventApp.NewHttpInfo(httpMethod, statusCode)
		containerTraffic.HttpInfoMap[httpMethod][statusCode] = httpInfo
	}
	httpInfo.HttpCount++
	now := time.Unix(0, msg.GetTimestamp())
	httpInfo.LastAcivity = &now
	httpInfo.LastResponseTime = responseTimeNano
}

// Format of HttpStartStop in PCF 1.6
// origin:"router__0" eventType:HttpStartStop timestamp:1481942228218776855 deployment:"cf"
// job:"router-partition-6c9fddda6d386d1b5b54" index:"0" ip:"172.28.1.57"
// httpStartStop:<startTimestamp:1481942228172587068 stopTimestamp:1481942228218776855
// requestId:<low:8306249200620409206  high:5348636437173287548 >  peerType:Client method:PUT
// uri:"api.system.laba.ecsteam.io/v2/spaces/2a7f2b63-e3f9-4e26-a73c-d3dd0be4b77f"
// remoteAddress:"172.28.1.51:46187" userAgent:"Mozilla/5.0" statusCode:201 contentLength:1300
// instanceId:"a455303f-144b-4936-9483-a61bfa23b35b" 11:"" >
func (ed *EventData) checkIfApiCall_PCF1_6(uri string) (bool, string) {
	if strings.HasPrefix(uri, ed.apiUrl) {
		// A PCF API has been called
		return true, strings.TrimPrefix(uri, ed.apiUrl)
	}
	return false, ""
}

// Format of HttpStartStop in PCF 1.7
// origin:"gorouter" eventType:HttpStartStop timestamp:1481950482175868678 deployment:"cf"
// job:"router-partition-72c346932f9a11cd262e" index:"0" ip:"172.28.3.57"
// httpStartStop:<startTimestamp:1481950482140382655 stopTimestamp:1481950482175838683
// requestId:<low:4704681661434642403 high:14399846285434694512 >  peerType:Client method:PUT
// uri:"http://73.169.24.191, 172.28.3.51, 172.29.0.2, 172.28.3.51/v2/spaces/605f2e92-a311-4bf8-a37d-296b0a692e25"
// remoteAddress:"172.28.3.51:38274" userAgent:"Mozilla/5.0" statusCode:201 contentLength:1311
// instanceId:"66a0e28b-81d7-44d6-46a4-ed46822b9b1f" 11:"" >
//
// NOTE: the uri has lost its hostname part which was the indication we are calling an api in PCF 1.6
// so we need to make some guesses that if we don't have an applicationId property but have an instanceId
// we must be calling an API.  Is this a good assumption?  Not a very clean solution.
func (ed *EventData) checkIfApiCall_PCF1_7(uri string) (bool, string) {
	toplog.Debug("Check if PCF 1.7 API call:%v", uri)
	parsedData := ed.urlPCF17Regexp.FindAllStringSubmatch(uri, -1)
	if len(parsedData) != 1 {
		toplog.Debug("Not a PCF 1.7 API that we care about")
		return false, ""
	}
	dataArray := parsedData[0]
	if len(dataArray) != 2 {
		toplog.Warn("checkIfApiCall_PCF1_7>>Unable to parse uri: %v", uri)
		return false, ""
	}
	apiUri := dataArray[1]
	toplog.Debug("This is a PCF 1.7 API call:%v", apiUri)
	return true, apiUri
}

func (ed *EventData) getContainerTraffic(appStats *eventApp.AppStats, instId string) *eventApp.TrafficStats {

	// Save the container data -- by instance id

	if appStats.ContainerTrafficMap == nil {
		appStats.ContainerTrafficMap = make(map[string]*eventApp.TrafficStats)
	}

	containerTraffic := appStats.ContainerTrafficMap[instId]
	if containerTraffic == nil {
		containerTraffic = eventApp.NewTrafficStats()
		appStats.ContainerTrafficMap[instId] = containerTraffic
		containerTraffic.ResponseL60Time = util.NewAvgTracker(time.Minute)
		containerTraffic.ResponseL10Time = util.NewAvgTracker(time.Second * 10)
		containerTraffic.ResponseL1Time = util.NewAvgTracker(time.Second)
	}

	return containerTraffic
}

// A PCF API has been called -- use this to trigger reload of metadata if appropriate
// Example: "/v2/spaces/59cde607-2cda-4e20-ab30-cc779c4026b0"
func (ed *EventData) pcfApiHasBeenCalled(msg *events.Envelope, apiUri string) {
	toplog.Debug("API called: %v", apiUri)

	parsedData := ed.apiUrlRegexp.FindAllStringSubmatch(apiUri, -1)
	if len(parsedData) != 1 {
		toplog.Debug("pcfApiHasBeenCalled>>Unable to parse (parsedData size) apiUri: %v", apiUri)
		return
	}
	dataArray := parsedData[0]
	if len(dataArray) != 3 {
		toplog.Debug("pcfApiHasBeenCalled>>Unable to parse (dataArray size) apiUri: %v", apiUri)
		return
	}
	dataType := dataArray[1]
	guid := dataArray[2]
	ed.pcfApiHasBeenCalledReloadMetadata(dataType, guid)

}

func (ed *EventData) pcfApiHasBeenCalledReloadMetadata(dataType, guid string) {
	toplog.Debug("Data type:%v GUID:%v", dataType, guid)
	switch dataType {
	case "spaces":
		// TODO reload metadata
		toplog.Debug("Reload SPACE metadata for space with GUID:%v", guid)
	case "organizations":
		// TODO reload metadata
		toplog.Debug("Reload ORG metadata for org with GUID:%v", guid)
	default:
	}
}

func (ed *EventData) handleRouteStats(msg *events.Envelope) {

	origin := msg.GetOrigin()
	if !strings.Contains(origin, "router") {
		// We are only interested in go-router events.
		// Other sources of HttpStartStop events:
		//   garden-windows
		//   etcd
		//	 routing_api  (PCF 1.8)
		return
	}

	httpEvent := msg.GetHttpStartStop()
	uri := httpEvent.GetUri()

	// Check if URI has a space in it
	if strings.IndexByte(uri, ' ') != -1 {
		// This is a bogus uri -- ignore it
		// PCF 1.7 messed up the HttpStartStop format using:
		// http://172.29.0.2, 172.28.3.51/oauth/token
		return
	}
	parsedData := ed.routeUrlRegexp.FindAllStringSubmatch(uri, -1)
	if len(parsedData) != 1 {
		toplog.Debug("handleRouteStats>>Unable to parse (parsedData size) apiUri: %v", uri)
		return
	}
	dataArray := parsedData[0]
	if len(dataArray) != 6 {
		toplog.Debug("handleRouteStats>>Unable to parse (dataArray size) apiUri: %v", uri)
		return
	}
	ipAddress := dataArray[1]
	host := dataArray[2]
	domain := dataArray[3]
	port := dataArray[4]
	path := dataArray[5]

	if ipAddress != "" {
		host = ipAddress
	}
	ed.updateRouteStats(domain, host, port, path, msg)

}

func (ed *EventData) GetAppRouteStats(uri string, domain string, host string, port string, path string, appId string) *eventRoute.AppRouteStats {

	domain = strings.ToLower(domain)
	host = strings.ToLower(host)

	domainStats := ed.DomainMap[domain]
	if domainStats == nil {
		toplog.Debug("domainStats not found. It will be dynamically added for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v]",
			uri, domain, host, port, path)
		if len(ed.DomainMap) > config.MaxDomainBucket {
			toplog.Warn("domainStats map at max size. The entry will NOT be added")
			return nil
		}
		domainGuid := util.Pseudo_uuid()
		domainStats = eventRoute.NewDomainStats(domainGuid)
		ed.DomainMap[domain] = domainStats
	}
	hostStats := domainStats.HostStatsMap[host]
	if hostStats == nil {
		// Check if we have a wildcard hostname
		hostStats = domainStats.HostStatsMap["*"]
		if hostStats != nil {
			toplog.Debug("hostStats wildcard found for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v]",
				uri, domain, host, port, path)
		} else {
			toplog.Debug("hostStats not found. It will be dynamically added for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v]",
				uri, domain, host, port, path)
			if len(domainStats.HostStatsMap) > config.MaxHostBucket {
				toplog.Warn("hostStats map at max size. The entry will NOT be added")
				return nil
			}
			// dynamically add new hosts/routes that we don't have pre-registered
			hostStats = eventRoute.NewHostStats(host)
			domainStats.HostStatsMap[host] = hostStats
		}
	}

	routeStats := hostStats.FindRouteStats(path)
	if routeStats == nil {
		toplog.Debug("routeStats not found. It will be dynamically added for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v]",
			uri, domain, host, port, path)
		// dynamically add root path
		routeStats = ed.eventProcessor.addInternalRoute(domain, host, "", 0)
		if routeStats == nil {
			return nil
		}
	}

	appRouteStats := routeStats.FindAppRouteStats(appId)
	if appRouteStats == nil {
		// QUESTION: Can a freshly deployed app (after cf top was running) end up here?
		// I think so -- the AppMetadata is loaded when a new app is deployed but the
		// route(s) for that app are never loaded/seeded (I need to verify this)
		toplog.Debug("appRouteStats not found. It will be dynamically added for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v]",
			uri, domain, host, port, path)
		appRouteStats = eventRoute.NewAppRouteStats(appId)
		routeStats.AppRouteStatsMap[appId] = appRouteStats
	}

	return appRouteStats
}

func (ed *EventData) updateRouteStats(domain string, host string, port string, path string, msg *events.Envelope) {

	httpEvent := msg.GetHttpStartStop()
	appUUID := httpEvent.GetApplicationId()
	appId := ""
	if appUUID != nil {
		appId = formatUUID(appUUID)
	}
	appRouteStats := ed.GetAppRouteStats(httpEvent.GetUri(), domain, host, port, path, appId)
	if appRouteStats == nil {
		// An internal error occurred -- we can't track this stat at this time
		return
	}

	httpMethod := httpEvent.GetMethod()
	httpMethodStats := appRouteStats.FindHttpMethodStats(httpMethod)
	if httpMethodStats == nil {
		toplog.Debug("httpMethodStats not found. It will be dynamically added for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v] method:[%v]",
			httpEvent.GetUri(), domain, host, port, path, httpMethod)
		httpMethodStats = eventRoute.NewHttpMethodStats(httpMethod)
		appRouteStats.HttpMethodStatsMap[httpMethod] = httpMethodStats
	}

	msgTime := time.Unix(0, msg.GetTimestamp())
	httpMethodStats.LastAccess = msgTime

	httpMethodStats.HttpStatusCode[httpEvent.GetStatusCode()] = httpMethodStats.HttpStatusCode[httpEvent.GetStatusCode()] + 1

	userAgentCount := appRouteStats.UserAgentMap[httpEvent.GetUserAgent()]
	if userAgentCount > 0 || (len(appRouteStats.UserAgentMap) < config.MaxUserAgentBucket) {
		appRouteStats.UserAgentMap[httpEvent.GetUserAgent()] = userAgentCount + 1
	}

	httpMethodStats.RequestCount = httpMethodStats.RequestCount + 1

	responseLength := httpEvent.GetContentLength()
	if responseLength > 0 {
		httpMethodStats.ResponseContentLength = httpMethodStats.ResponseContentLength + httpEvent.GetContentLength()
	}
	/*
		toplog.Debug("Updated stats for uri:[%v] domain:[%v] host:[%v] port:[%v] path:[%v] method:[%v]",
			httpEvent.GetUri(), domain, host, port, path, httpMethod)
	*/
}
