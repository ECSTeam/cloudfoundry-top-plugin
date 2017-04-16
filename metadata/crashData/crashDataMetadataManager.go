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
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

var (
	LoadEventsUntilTime *time.Time

	crashDataMetadataCache []EventData
	// Map: [AppGuid][instanceIndex] = array of crash timestamps
	crashDataByAppId map[string]map[int][]*time.Time
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

func FindByApp(appGuid string) []*time.Time {
	var allCrashTimestamps []*time.Time
	if crashDataByAppId != nil {
		crashMap := crashDataByAppId[appGuid]
		if crashMap != nil {
			for _, crashTimestamps := range crashMap {
				if crashTimestamps != nil {
					allCrashTimestamps = append(allCrashTimestamps, crashTimestamps...)
				}
			}
		}
	}
	sort.Sort(util.TimeSlice(allCrashTimestamps))
	return allCrashTimestamps
}

func FindSinceByApp(appGuid string, since time.Duration) []*time.Time {
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

func FindByAppAndInstance(appGuid string, instanceIndex int) []*time.Time {
	if crashDataByAppId != nil {
		crashMap := crashDataByAppId[appGuid]
		if crashMap != nil {
			crashTimestamps := crashMap[instanceIndex]
			if crashTimestamps != nil {
				return crashTimestamps
			}
		}
	}
	return nil
}

func FindLastCrashByAppAndInstance(appGuid string, instanceIndex int) *time.Time {
	crashTimestamps := FindByAppAndInstance(appGuid, instanceIndex)
	if crashTimestamps != nil && len(crashTimestamps) > 0 {
		return crashTimestamps[len(crashTimestamps)-1]
	}
	return nil
}

func FindSinceByAppAndInstance(appGuid string, instanceIndex int, since time.Duration) []*time.Time {
	crashTimestamps := FindByAppAndInstance(appGuid, instanceIndex)
	return filterSince(crashTimestamps, since)
}

func FindCountSinceByAppAndInstance(appGuid string, instanceIndex int, since time.Duration) int {
	crashTimestamps := FindSinceByAppAndInstance(appGuid, instanceIndex, since)
	if crashTimestamps != nil {
		return len(crashTimestamps)
	}
	return 0
}

func filterSince(crashTimestamps []*time.Time, since time.Duration) []*time.Time {

	if crashTimestamps != nil {
		sinceTime := time.Now().Add(since)
		crashTimestampsSize := len(crashTimestamps)
		crashTimestampsSince := make([]*time.Time, 0, crashTimestampsSize)
		for i, _ := range crashTimestamps {
			// Reverse loop through array
			crashTS := crashTimestamps[crashTimestampsSize-i-1]
			if crashTS == nil || crashTS.Before(sinceTime) {
				break
			}
			//crashCount = i + 1
			crashTimestampsSince = append(crashTimestampsSince, crashTS)
		}
		return crashTimestampsSince
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

	crashDataByAppId = make(map[string]map[int][]*time.Time)

	layout := "2006-01-02T15:04:05Z"
	for _, crashData := range data {
		crashMap := crashDataByAppId[crashData.Actor]
		if crashMap == nil {
			crashMap = make(map[int][]*time.Time)
			crashDataByAppId[crashData.Actor] = crashMap
		}
		crashTimestamp, err := time.Parse(layout, crashData.Timestamp)
		if err != nil {
			fmt.Println(err)
		}

		instanceIndex := crashData.Metadata.Index
		crashTimestamps := crashMap[instanceIndex]
		if crashTimestamps == nil {
			crashTimestamps = make([]*time.Time, 0, 20)
		}
		crashDataByAppId[crashData.Actor][instanceIndex] = append(crashTimestamps, &crashTimestamp)
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
