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

package appStatistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

// ****************************************************************
// The following are used used calling API: /v2/apps/APP_GUID/stats
// ****************************************************************
type AppInstanceStatistics struct {
	CacheTime *time.Time
	Data      map[string]*AppInstanceStatistic
}
type AppInstanceStatistic struct {
	State string `json:"state"`
	Stats struct {
		Name                string     `json:"name"`
		URIs                []string   `json:"uris"`
		Host                string     `json:"host"`
		Port                int        `json:"port"`
		Uptime              int64      `json:"uptime"`
		StartTime           *time.Time // This will be populated on post-processing of response
		MemoryQuota         int64      `json:"mem_quota"`
		DiskQuota           int64      `json:"disk_quota"`
		FiledescriptorQuota int        `json:"fds_quota"`
		Usage               struct {
			Time   string  `json:"time"`
			CPU    float64 `json:"cpu"`
			Memory int64   `json:"mem"`
			Disk   int64   `json:"disk"`
		} `json:"usage"`
	} `json:"stats"`
}

var (
	// A map of AppIds
	appStatisticsMetadataCache = make(map[string]*AppInstanceStatistics)
	mu                         sync.Mutex
)

func FindAppStatisticMetadata(appId string) map[string]*AppInstanceStatistic {

	stats := FindAppStatisticMetadataInternal(appId)
	if stats != nil {
		appInstanceStatistics := stats.Data
		if appInstanceStatistics != nil {
			return appInstanceStatistics
		}
	}
	//return make(map[string]*AppInstanceStatistic)
	return nil
}

func FindAppStatisticMetadataInternal(appId string) *AppInstanceStatistics {
	mu.Lock()
	defer mu.Unlock()
	return appStatisticsMetadataCache[appId]
}

func LoadAppStatisticCache(cliConnection plugin.CliConnection, appId string) error {

	now := time.Now()
	data, err := getAppStatisticMetadata(cliConnection, appId)
	if err != nil {
		toplog.Warn("*** app instance metadata error: %v", err.Error())
		return err
	}

	instStats := &AppInstanceStatistics{CacheTime: &now, Data: data}
	mu.Lock()
	defer mu.Unlock()
	appStatisticsMetadataCache[appId] = instStats
	return nil
}

func getAppStatisticMetadata(cliConnection plugin.CliConnection, appId string) (map[string]*AppInstanceStatistic, error) {

	url := "/v2/apps/" + appId + "/stats"

	output, err := common.CallAPI(cliConnection, url)
	if err != nil {
		return nil, err
	}

	if strings.Contains(output, "error_code") {
		if strings.Contains(output, "CF-AppStoppedStatsError") {
			// This error is OK
			return make(map[string]*AppInstanceStatistic), nil
		} else {
			errMsg := fmt.Sprintf("Error from API call: %v", output)
			return nil, errors.New(errMsg)
		}
	}

	response := make(map[string]*AppInstanceStatistic)
	outputBytes := []byte(output)
	err = json.Unmarshal(outputBytes, &response)
	if err != nil {
		toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
		return response, err
	}

	// Set the startTime relative to now and uptime of the container
	now := time.Now().Truncate(time.Second)
	for _, stat := range response {
		uptimeSeconds := stat.Stats.Uptime
		startTime := now.Add(time.Duration(-uptimeSeconds) * time.Second)
		stat.Stats.StartTime = &startTime
	}

	return response, nil

}
