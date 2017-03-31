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

package app

import (
	"encoding/json"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
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

func (mdMgr *AppMetadataManager) GetAppMetadataMap() map[string]*AppMetadata {
	return mdMgr.appMetadataMap
}

func (mdMgr *AppMetadataManager) AllApps() []*AppMetadata {
	appsMetadataArray := []*AppMetadata{}
	for _, appMetadata := range mdMgr.appMetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata)
	}
	return appsMetadataArray
}

func (mdMgr *AppMetadataManager) FindAppMetadata(appId string) *AppMetadata {
	return mdMgr.FindAppMetadataInternal(appId, true)
}

func (mdMgr *AppMetadataManager) FindAppMetadataInternal(appId string, requestLoadIfNotFound bool) *AppMetadata {
	appMetadata := mdMgr.appMetadataMap[appId]
	if appMetadata == nil {
		appMetadata = NewAppMetadataById(appId)
		if requestLoadIfNotFound {
			// TODO: Queue metadata load for this id
		} else {
			// We mark this metadata as 60 mins old
			appMetadata.CacheTime = appMetadata.CacheTime.Add(-60 * time.Minute)
		}
	}
	return appMetadata
}

func (mdMgr *AppMetadataManager) LoadAppCache(cliConnection plugin.CliConnection) {
	appMetadataArray, err := mdMgr.getAppsMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** app metadata error: %v", err.Error())
		return
	}

	metadataMap := make(map[string]*AppMetadata)
	for _, appMetadata := range appMetadataArray {
		//toplog.Debug("From Map - app id: %v name:%v", appMetadata.Guid, appMetadata.Name)
		metadataMap[appMetadata.Guid] = appMetadata
	}

	mdMgr.appMetadataMap = metadataMap
}

func (mdMgr *AppMetadataManager) GetAppMetadataInternal(cliConnection plugin.CliConnection, appId string) (*AppMetadata, error) {
	url := "/v2/apps/" + appId
	emptyApp := NewAppMetadataById(appId)

	outputStr, err := common.CallAPI(cliConnection, url)
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
	return GetAppsMetadataFromUrl(cliConnection, "/v2/apps")
}

func GetAppsMetadataFromUrl(cliConnection plugin.CliConnection, url string) ([]*AppMetadata, error) {

	appsMetadataArray := []*AppMetadata{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var appResp AppResponse
		err = json.Unmarshal(outputBytes, &appResp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return appsMetadataArray, "", err
		}
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			appMetadata := NewAppMetadata(app.Entity)
			appsMetadataArray = append(appsMetadataArray, appMetadata)
		}
		return appResp, appResp.NextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return appsMetadataArray, err

}
