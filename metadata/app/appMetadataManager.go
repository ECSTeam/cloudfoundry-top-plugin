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
	"fmt"
	"sync"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/loader"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/mdGlobalManagerInterface"
)

type AppMetadataManager struct {
	*common.CommonV2ResponseManager
	mu sync.Mutex
}

func NewAppMetadataManager(mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface) *AppMetadataManager {
	url := "/v2/apps"
	mdMgr := &AppMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, url, mdMgr, false)
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
	return mdMgr.FindItemInternal(appId, false, true).(*AppMetadata)
}

func (mdMgr *AppMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewAppMetadataById(guid)
}

func (mdMgr *AppMetadataManager) CreateResponseObject() common.IResponse {
	return &AppResponse{}
}

func (mdMgr *AppMetadataManager) CreateResourceObject() common.IResource {
	return &AppResource{}
}

func (mdMgr *AppMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewAppMetadataById(guid)
}

func (mdMgr *AppMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*AppResponse)
	for _, item := range resp.Resources {
		item.Entity.Guid = item.Meta.Guid
		metadata := NewAppMetadata(item.Entity)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *AppMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*AppResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewAppMetadata(resourceType.Entity)
	return metadata
}

func (mdMgr *AppMetadataManager) CreateTestData(dataSize int) {
	metadataMap := make(map[string]common.IMetadata)
	for i := 0; i < dataSize; i++ {
		guid := fmt.Sprintf("GUID-%02v", i)
		app := &App{EntityCommon: common.EntityCommon{Guid: guid}, Name: guid}
		appMetadata := NewAppMetadata(*app)
		metadataMap[guid] = appMetadata
	}
	mdMgr.MetadataMap = metadataMap
}
