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

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type AppMetadataManager struct {
	appMetadataMap map[string]*AppMetadata
	mu             sync.Mutex
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
	return mdMgr.findAppMetadataInternal(appId, true)
}

func (mdMgr *AppMetadataManager) findAppMetadataInternal(appId string, requestLoadIfNotFound bool) *AppMetadata {
	appMetadata := mdMgr.appMetadataMap[appId]
	if appMetadata == nil {
		appMetadata = NewAppMetadataById(appId)
		if requestLoadIfNotFound {
			// TODO: Queue metadata load for this id
		} else {
			// We mark this metadata as 60 mins old
			appMetadata.cacheTime = appMetadata.cacheTime.Add(-60 * time.Minute)
		}
	}
	return appMetadata
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

	err := callPagableAPI(cliConnection, url, handleRequest)

	return appsMetadataArray, err

}
