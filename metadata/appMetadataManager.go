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

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type AppMetadataManager struct {
	appMetadataMap            map[string]*AppMetadata
	totalMemoryAllStartedApps float64
	totalDiskAllStartedApps   float64

	mu sync.Mutex
}

func NewAppMetadataManager() *AppMetadataManager {

	mgr := &AppMetadataManager{}
	mgr.appMetadataMap = make(map[string]*AppMetadata)

	return mgr
}

func (mdMgr *AppMetadataManager) AppMetadataSize() int {
	return len(mdMgr.appMetadataMap)
}

func (mdMgr *AppMetadataManager) AllApps() []*AppMetadata {
	appsMetadataArray := []*AppMetadata{}
	for _, appMetadata := range mdMgr.appMetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata)
	}
	return appsMetadataArray
}

func (mdMgr *AppMetadataManager) FindAppMetadata(appId string) *AppMetadata {

	appMetadata := mdMgr.appMetadataMap[appId]
	if appMetadata == nil {
		appMetadata = NewAppMetadataById(appId)
	}
	return appMetadata
}

func (mdMgr *AppMetadataManager) GetTotalMemoryAllStartedApps() float64 {
	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	//toplog.Debug("entering GetTotalMemoryAllStartedApps")
	if mdMgr.totalMemoryAllStartedApps == 0 {
		total := float64(0)
		for _, app := range mdMgr.appMetadataMap {
			if app.State == "STARTED" {
				total = total + ((app.MemoryMB * MEGABYTE) * app.Instances)
			}
		}
		mdMgr.totalMemoryAllStartedApps = total
	}
	//toplog.Debug("leaving GetTotalMemoryAllStartedApps")
	return mdMgr.totalMemoryAllStartedApps
}

func (mdMgr *AppMetadataManager) GetTotalDiskAllStartedApps() float64 {
	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	//toplog.Debug("entering GetTotalDiskAllStartedApps")
	if mdMgr.totalDiskAllStartedApps == 0 {
		total := float64(0)
		for _, app := range mdMgr.appMetadataMap {
			if app.State == "STARTED" {
				total = total + ((app.DiskQuotaMB * MEGABYTE) * app.Instances)
			}
		}
		mdMgr.totalDiskAllStartedApps = total
	}
	//toplog.Debug("leaving GetTotalDiskAllStartedApps")
	return mdMgr.totalDiskAllStartedApps
}

func (mdMgr *AppMetadataManager) LoadAppCache(cliConnection plugin.CliConnection) {
	appMetadataArray, err := mdMgr.getAppsMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** app metadata error: %v", err.Error()))
		return
	}

	metadataMap := make(map[string]*AppMetadata)
	for _, appMetadata := range appMetadataArray {
		//toplog.Debug(fmt.Sprintf("From Map - app id: %v name:%v", appMetadata.Guid, appMetadata.Name))
		metadataMap[appMetadata.Guid] = appMetadata
	}

	mdMgr.appMetadataMap = metadataMap
}

func (mdMgr *AppMetadataManager) getAppMetadata(cliConnection plugin.CliConnection, appId string) (*AppMetadata, error) {
	url := "/v2/apps/" + appId
	emptyApp := NewAppMetadataById(appId)

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
	mdMgr.flushCounters()
	appMetadata := NewAppMetadata(appResource.Entity)
	return appMetadata, nil
}

func (mdMgr *AppMetadataManager) getAppsMetadata(cliConnection plugin.CliConnection) ([]*AppMetadata, error) {

	url := "/v2/apps"
	appsMetadataArray := []*AppMetadata{}

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
			appsMetadataArray = append(appsMetadataArray, appMetadata)
		}
		return appResp, nil
	}

	callPagableAPI(cliConnection, url, handleRequest)

	mdMgr.flushCounters()
	return appsMetadataArray, nil

}

func (mdMgr *AppMetadataManager) flushCounters() {
	// Flush the total counters
	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	mdMgr.totalMemoryAllStartedApps = 0
	mdMgr.totalDiskAllStartedApps = 0
}
