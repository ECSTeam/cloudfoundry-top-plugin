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

package appInstances

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
// The following are used used calling API: /v2/apps/APP_GUID/instances
// ****************************************************************
type AppInstances struct {
	CacheTime *time.Time
	Data      map[string]*AppInstance
}
type AppInstance struct {
	Details   string     `json:"details"`
	Since     float64    `json:"since"`
	State     string     `json:"state"`
	Uptime    int64      `json:"uptime"`
	StartTime *time.Time // This will be populated on post-processing of response

}

var (
	// A map of AppIds
	appInstancesMetadataCache = make(map[string]*AppInstances)
	mu                        sync.Mutex
)

func FindAppInstancesMetadata(appId string) *AppInstances {
	return FindAppInstancesMetadataInternal(appId)
}

func ClearAppInstancesMetadata(appId string) {
	mu.Lock()
	defer mu.Unlock()
	appInstancesMetadataCache[appId] = nil
}

func FindAppInstancesMetadataInternal(appId string) *AppInstances {
	mu.Lock()
	defer mu.Unlock()
	return appInstancesMetadataCache[appId]
}

func LoadAppInstancesCache(cliConnection plugin.CliConnection, appId string) error {

	now := time.Now()
	data, err := getAppInstancesMetadata(cliConnection, appId)
	if err != nil {
		toplog.Warn("*** app instance metadata error: %v  response: %v   appId: %v", err.Error(), data, appId)
		return err
	}

	instStats := &AppInstances{CacheTime: &now, Data: data}
	mu.Lock()
	defer mu.Unlock()
	appInstancesMetadataCache[appId] = instStats
	return nil
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	appInstancesMetadataCache = make(map[string]*AppInstances)
}

func getAppInstancesMetadata(cliConnection plugin.CliConnection, appId string) (map[string]*AppInstance, error) {

	url := "/v2/apps/" + appId + "/instances"

	output, err := common.CallAPI(cliConnection, url)
	if err != nil {
		return nil, err
	}

	if strings.Contains(output, "error_code") {
		// "Instances error: Request failed for app: cf-nodejs as the app is in stopped state."
		if strings.Contains(output, "220001") {
			// This error is OK
			return make(map[string]*AppInstance), nil
		} else {
			errMsg := fmt.Sprintf("Error from API call: %v", output)
			return nil, errors.New(errMsg)
		}
	}

	response := make(map[string]*AppInstance)
	outputBytes := []byte(output)
	err = json.Unmarshal(outputBytes, &response)
	if err != nil {
		toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
		return response, err
	}

	// Set the startTime relative to now and uptime of the container
	for _, stat := range response {
		// Ignore "uptime" field if container is in state DOWN
		if stat.State == "DOWN" {
			stat.Uptime = 0
		} else {
			startTime := time.Unix(int64(stat.Since), 0)
			stat.StartTime = &startTime
		}
	}

	return response, nil

}
