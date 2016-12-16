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

package metadata

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

const MEGABYTE = (1024 * 1024)

type AppMetadata struct {
	App
	cacheTime time.Time
}

func NewAppMetadata(appStats App) *AppMetadata {
	appMetadata := &AppMetadata{}
	appMetadata.App = appStats
	appMetadata.cacheTime = time.Now()
	return appMetadata
}

type AppResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []AppResource `json:"resources"`
}

type AppResource struct {
	Meta   Meta `json:"metadata"`
	Entity App  `json:"entity"`
}

type App struct {
	Guid      string `json:"guid"`
	Name      string `json:"name,omitempty"`
	SpaceGuid string `json:"space_guid,omitempty"`
	SpaceName string
	OrgGuid   string
	OrgName   string

	StackGuid   string  `json:"stack_guid,omitempty"`
	MemoryMB    float64 `json:"memory,omitempty"`
	DiskQuotaMB float64 `json:"disk_quota,omitempty"`

	Environment map[string]interface{} `json:"environment_json,omitempty"`
	Instances   float64                `json:"instances,omitempty"`
	State       string                 `json:"state,omitempty"`
	EnableSsh   bool                   `json:"enable_ssh,omitempty"`

	PackageState        string `json:"package_state,omitempty"`
	StagingFailedReason string `json:"staging_failed_reason,omitempty"`
	StagingFailedDesc   string `json:"staging_failed_description,omitempty"`
	DetectedStartCmd    string `json:"detected_start_command,omitempty"`
	//DockerCredentials string  `json:"docker_credentials_json,omitempty"`
	//audit.app.create event fields
	Console           bool   `json:"console,omitempty"`
	Buildpack         string `json:"buildpack,omitempty"`
	DetectedBuildpack string `json:"detected_buildpack,omitempty"`

	HealthcheckType    string  `json:"health_check_type,omitempty"`
	HealthcheckTimeout float64 `json:"health_check_timeout,omitempty"`
	Production         bool    `json:"production,omitempty"`
	//app.crash event fields
	//Index           float64 `json:"index,omitempty"`
	//ExitStatus      string  `json:"exit_status,omitempty"`
	//ExitDescription string  `json:"exit_description,omitempty"`
	//ExitReason      string  `json:"reason,omitempty"`
	// "package_updated_at": "2016-11-15T19:56:52Z",
	PackageUpdatedAt string `json:"package_updated_at"`
}

var (
	appMetadataMap            map[string]AppMetadata
	totalMemoryAllStartedApps float64
	totalDiskAllStartedApps   float64
	mu                        sync.Mutex
	refreshNow                chan bool
	refreshQueue              map[string]string
	cliConnection             plugin.CliConnection
)

// TODO: Convert this entire thing into a class
func init() {
	refreshQueue = make(map[string]string)
	refreshNow = make(chan bool)
	appMetadataMap = make(map[string]AppMetadata)
}

func SetConnection(conn plugin.CliConnection) {
	cliConnection = conn
	go loadMetadataThread()
}

func AppMetadataSize() int {
	return len(appMetadataMap)
}

func AllApps() []AppMetadata {
	appsMetadataArray := []AppMetadata{}
	for _, appMetadata := range appMetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata)
	}
	return appsMetadataArray
}

func FindAppMetadata(appId string) AppMetadata {
	appMetadata := appMetadataMap[appId]
	if appMetadataMap == nil {
		appMetadata = AppMetadata{}
	}
	return appMetadata
}

// Request a refresh of specific app metadata
func RequestRefreshAppMetadata(appId string) {
	refreshQueue[appId] = appId
	wakeRefreshThread()
}

func wakeRefreshThread() {
	refreshNow <- true
}

func loadMetadataThread() {

	minimumLoadTimeMS := time.Millisecond * 10000
	veryLongtime := time.Hour * 10000
	minNextLoadTime := veryLongtime

	for {

		toplog.Info(fmt.Sprintf("Metadata - sleep time: %v", minNextLoadTime))

		select {
		case <-refreshNow:
			//mui.updateDisplay(g)
		case <-time.After(minNextLoadTime):
			//mui.updateDisplay(g)
		}

		minNextLoadTime = veryLongtime
		toplog.Warn("hello from cache thread")
		for _, appId := range refreshQueue {
			appMetadata := appMetadataMap[appId]
			timeSinceLastLoad := time.Now().Sub(appMetadata.cacheTime)
			toplog.Debug(fmt.Sprintf("Metadata - appId: %v - inqueue check time since last load: %v", appId, timeSinceLastLoad))
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug(fmt.Sprintf("Metadata - appId: %v - Needs to be loaded now", appId))
				newAppMetadata, err := getAppMetadata(cliConnection, appId)
				if err != nil {
					toplog.Warn(fmt.Sprintf("Metadata - appId: %v - Error: %v", appId, err))
				} else {
					toplog.Info(fmt.Sprintf("Metadata - appId: %v - Load start", appId))
					appMetadataMap[appId] = newAppMetadata
					toplog.Info(fmt.Sprintf("Metadata - appId: %v - Load complete", appId))
				}
			} else {
				toplog.Debug(fmt.Sprintf("Metadata - appId %v - Too soon to reload", appId))
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug(fmt.Sprintf("Metadata - appId %v - Try to load in: %v", appId, nextLoadTime))
				if minNextLoadTime > nextLoadTime {
					toplog.Debug(fmt.Sprintf("Metadata - appId %v - value was min: %v", appId, nextLoadTime))
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}

func GetTotalMemoryAllStartedApps() float64 {
	mu.Lock()
	defer mu.Unlock()
	//toplog.Debug("entering GetTotalMemoryAllStartedApps")
	if totalMemoryAllStartedApps == 0 {
		total := float64(0)
		for _, app := range appMetadataMap {
			if app.State == "STARTED" {
				total = total + ((app.MemoryMB * MEGABYTE) * app.Instances)
			}
		}
		totalMemoryAllStartedApps = total
	}
	//toplog.Debug("leaving GetTotalMemoryAllStartedApps")
	return totalMemoryAllStartedApps
}

func GetTotalDiskAllStartedApps() float64 {
	mu.Lock()
	defer mu.Unlock()
	//toplog.Debug("entering GetTotalDiskAllStartedApps")
	if totalDiskAllStartedApps == 0 {
		total := float64(0)
		for _, app := range appMetadataMap {
			if app.State == "STARTED" {
				total = total + ((app.DiskQuotaMB * MEGABYTE) * app.Instances)
			}
		}
		totalDiskAllStartedApps = total
	}
	//toplog.Debug("leaving GetTotalDiskAllStartedApps")
	return totalDiskAllStartedApps
}

func LoadAppCache(cliConnection plugin.CliConnection) {
	appMetadataArray, err := getAppsMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** app metadata error: %v", err.Error()))
		return
	}

	metadataMap := make(map[string]AppMetadata)
	for _, appMetadata := range appMetadataArray {
		metadataMap[appMetadata.Guid] = appMetadata
	}
	appMetadataMap = metadataMap
}

func getAppMetadata(cliConnection plugin.CliConnection, appId string) (AppMetadata, error) {
	url := "/v2/apps/" + appId
	emptyApp := AppMetadata{}

	outputStr, err := callAPI(cliConnection, url)
	if err != nil {
		return emptyApp, err
	}
	outputBytes := []byte(outputStr)
	var appResource AppResource
	err = json.Unmarshal(outputBytes, &appResource)
	if err != nil {
		return emptyApp, err
	}
	appResource.Entity.Guid = appResource.Meta.Guid
	flushCounters()
	appMetadata := NewAppMetadata(appResource.Entity)
	return *appMetadata, nil
}

func getAppsMetadata(cliConnection plugin.CliConnection) ([]AppMetadata, error) {

	url := "/v2/apps"
	appsMetadataArray := []AppMetadata{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var appResp AppResponse
		err := json.Unmarshal(outputBytes, &appResp)
		if err != nil {
			toplog.Warn(fmt.Sprintf("*** %v unmarshal parsing output: %v", url, string(outputBytes[:])))
			return appsMetadataArray, err
		}
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			appMetadata := NewAppMetadata(app.Entity)
			appsMetadataArray = append(appsMetadataArray, *appMetadata)
		}
		return appResp, nil
	}

	callPagableAPI(cliConnection, url, handleRequest)

	flushCounters()
	return appsMetadataArray, nil

}

func flushCounters() {
	// Flush the total counters
	mu.Lock()
	defer mu.Unlock()
	totalMemoryAllStartedApps = 0
	totalDiskAllStartedApps = 0
}
