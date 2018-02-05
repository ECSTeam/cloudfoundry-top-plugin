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

package space

import (
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/isolationSegment"
)

type SpaceMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewSpaceMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *SpaceMetadataManager {
	url := "/v2/spaces"
	mdMgr := &SpaceMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, common.SPACE, url, mdMgr, false)
	return mdMgr
}

func (mdMgr *SpaceMetadataManager) FindItem(guid string) *SpaceMetadata {
	return mdMgr.FindItemInternal(guid, false, true).(*SpaceMetadata)
}

func (mdMgr *SpaceMetadataManager) GetAll() []*SpaceMetadata {
	mdMgr.MetadataMapMutex.Lock()
	defer mdMgr.MetadataMapMutex.Unlock()
	appsMetadataArray := []*SpaceMetadata{}
	for _, appMetadata := range mdMgr.MetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata.(*SpaceMetadata))
	}
	return appsMetadataArray
}

func (mdMgr *SpaceMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewSpaceMetadataById(guid)
}

func (mdMgr *SpaceMetadataManager) CreateResponseObject() common.IResponse {
	return &SpaceResponse{}
}

func (mdMgr *SpaceMetadataManager) CreateResourceObject() common.IResource {
	return &SpaceResource{}
}

func (mdMgr *SpaceMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewSpaceMetadataById(guid)
}

func (mdMgr *SpaceMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*SpaceResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *SpaceMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*SpaceResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	if resourceType.Entity.IsolationSegmentGuid == "" {
		resourceType.Entity.IsolationSegmentGuid = isolationSegment.DefaultIsolationSegmentGuid
	}
	metadata := NewSpaceMetadata(resourceType.Entity)
	return metadata
}
