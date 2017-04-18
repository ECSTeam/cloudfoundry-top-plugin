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

package crashData

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

var (
	LoadEventsUntilTime *time.Time

	crashDataMetadataCache []EventData
	// Map: [AppGuid] = array of crash timestamps
	crashDataByAppId map[string][]*ContainerCrashInfo
)

func All() []EventData {
	return crashDataMetadataCache
}

func Find(guid string) EventData {
	for _, crashData := range crashDataMetadataCache {
		if crashData.Guid == guid {
			return crashData
		}
	}
	return EventData{EntityCommon: common.EntityCommon{Guid: guid}}
}

func FindByApp(appGuid string) []*ContainerCrashInfo {
	crashInfo := crashDataByAppId[appGuid]
	sort.Sort(ContainerCrashInfoSlice(crashInfo))
	return crashInfo
}

func FindSinceByApp(appGuid string, since time.Duration) []*ContainerCrashInfo {
	crashTimestamps := FindByApp(appGuid)
	return filterSince(crashTimestamps, since)
}

func FindCountSinceByApp(appGuid string, since time.Duration) int {
	crashTimestamps := FindSinceByApp(appGuid, since)
	if crashTimestamps != nil {
		return len(crashTimestamps)
	}
	return 0
}

func FindLastCrashByApp(appGuid string) *ContainerCrashInfo {
	crashInfoList := FindByApp(appGuid)
	if crashInfoList != nil && len(crashInfoList) > 0 {
		return crashInfoList[len(crashInfoList)-1]
	}
	return nil
}

func filterSince(crashInfoList []*ContainerCrashInfo, since time.Duration) []*ContainerCrashInfo {

	if crashInfoList != nil {
		sinceTime := time.Now().Add(since)
		crashInfoListSize := len(crashInfoList)
		crashInfoListSince := make([]*ContainerCrashInfo, 0, crashInfoListSize)
		for i, _ := range crashInfoList {
			// Reverse loop through array
			crashInfo := crashInfoList[crashInfoListSize-i-1]
			if crashInfo == nil || crashInfo.CrashTime.Before(sinceTime) {
				break
			}
			//crashCount = i + 1
			crashInfoListSince = append(crashInfoListSince, crashInfo)
		}
		return crashInfoListSince
	}
	return nil
}

func LoadCrashDataCache(cliConnection plugin.CliConnection) {
	data, err := getCrashDataMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** CrashData metadata error: %v", err.Error())
		return
	}
	crashDataMetadataCache = data

	crashDataByAppId = make(map[string][]*ContainerCrashInfo)

	layout := "2006-01-02T15:04:05Z"
	for _, crashData := range data {
		crashInfoList := crashDataByAppId[crashData.Actor]
		if crashInfoList == nil {
			crashInfoList = make([]*ContainerCrashInfo, 0)
			crashDataByAppId[crashData.Actor] = crashInfoList
		}
		crashTimestamp, err := time.Parse(layout, crashData.Timestamp)
		if err != nil {
			fmt.Println(err)
		}
		instanceIndex := crashData.Metadata.Index
		exitDescription := crashData.Metadata.Exit_description
		crashInfo := NewContainerCrashInfo(instanceIndex, &crashTimestamp, exitDescription)
		crashDataByAppId[crashData.Actor] = append(crashDataByAppId[crashData.Actor], crashInfo)
	}
}

func getCrashDataMetadata(cliConnection plugin.CliConnection) ([]EventData, error) {
	timestampFormat := "2006-01-02 15:04:05-07:00"
	urlPath := "/v2/events?q=type:app.crash&q=timestamp%%3E=%v&q=timestamp%%3C=%v"

	// TODO: "now" should be the timestamp of when top started
	// so we don't end up with dumplications of crash data
	eventsUtilTime := LoadEventsUntilTime
	eventsUtilTimeStr := eventsUtilTime.Format(timestampFormat)
	eventsUtilTimeStrEncoded := url.PathEscape(eventsUtilTimeStr)

	oneDayAgo := time.Now().Add(-24 * time.Hour)
	oneDayAgoStr := oneDayAgo.Format(timestampFormat)
	oneDayAgoStrEncoded := url.PathEscape(oneDayAgoStr)
	urlPath = fmt.Sprintf(urlPath, oneDayAgoStrEncoded, eventsUtilTimeStrEncoded)

	metadata := []EventData{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response EventDataResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", urlPath, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			metadata = append(metadata, item.Entity)
		}

		nextUrl = response.NextUrl
		// There is bug in value return for nextURL (tested in PCF 1.10.3) where the URL contains
		// "order-by" fields which the API rejects for subsequent pages
		nextUrl = strings.Replace(nextUrl, "order-by=timestamp&", "", -1)
		nextUrl = strings.Replace(nextUrl, "order-by=id&", "", -1)
		return response, nextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, urlPath, handleRequest)

	toplog.Debug("Total crash events loaded: %v", len(metadata))
	return metadata, err

}
