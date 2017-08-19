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
	"fmt"
	"sync"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/loader"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type AppMetadataManager struct {
	*common.BaseMetadataManager
	mu sync.Mutex
}

func NewAppMetadataManager() *AppMetadataManager {

	mdMgr := &AppMetadataManager{}
	mdMgr.BaseMetadataManager = common.NewBaseMetadataManager(mdMgr)
	loader.RegisterMetadataHandler(loader.APP, mdMgr)
	return mdMgr
}

func (mdMgr *AppMetadataManager) AllApps() []*AppMetadata {
	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	appsMetadataArray := []*AppMetadata{}
	for _, appMetadata := range mdMgr.MetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata.(*AppMetadata))
	}
	return appsMetadataArray
}

func (mdMgr *AppMetadataManager) FindItem(appId string) *AppMetadata {
	return mdMgr.FindItemInternal(appId, false).(*AppMetadata)
}

func (mdMgr *AppMetadataManager) NewItemById(guid string) common.BaseMetadataItemI {
	return NewAppMetadataById(guid)
}

func (mdMgr *AppMetadataManager) GetItemInternal(cliConnection plugin.CliConnection, guid string) (common.BaseMetadataItemI, error) {
	url := "/v2/apps/" + guid
	emptyApp := mdMgr.NewItemById(guid)

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

func (mdMgr *AppMetadataManager) LoadInternal(cliConnection plugin.CliConnection) ([]common.BaseMetadataItemI, error) {
	return GetAppsMetadataFromUrl(cliConnection, "/v2/apps")
}

func GetAppsMetadataFromUrl(cliConnection plugin.CliConnection, url string) ([]common.BaseMetadataItemI, error) {

	metadataItemArray := []common.BaseMetadataItemI{}
	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var appResp AppResponse
		err = json.Unmarshal(outputBytes, &appResp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadataItemArray, "", err
		}
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			appMetadata := NewAppMetadata(app.Entity)
			metadataItemArray = append(metadataItemArray, appMetadata)
		}
		return appResp, appResp.NextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)
	return metadataItemArray, err

}

func (mdMgr *AppMetadataManager) CreateTestData(dataSize int) {
	metadataMap := make(map[string]common.BaseMetadataItemI)
	for i := 0; i < dataSize; i++ {
		guid := fmt.Sprintf("GUID-%02v", i)
		app := &App{Guid: guid, Name: guid}
		appMetadata := NewAppMetadata(*app)
		metadataMap[guid] = appMetadata
	}
	mdMgr.MetadataMap = metadataMap
}
