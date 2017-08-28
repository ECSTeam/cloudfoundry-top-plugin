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

package spaceQuota

import (
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/mdGlobalManagerInterface"
)

type SpaceQuotaMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewSpaceQuotaMetadataManager(mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface) *SpaceQuotaMetadataManager {
	url := "/v2/space_quota_definitions"
	mdMgr := &SpaceQuotaMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, url, mdMgr, true)
	return mdMgr
}

func (mdMgr *SpaceQuotaMetadataManager) FindItem(appId string) *SpaceQuotaMetadata {
	return mdMgr.FindItemInternal(appId, false, true).(*SpaceQuotaMetadata)
}

func (mdMgr *SpaceQuotaMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewSpaceQuotaMetadataById(guid)
}

func (mdMgr *SpaceQuotaMetadataManager) CreateResponseObject() common.IResponse {
	return &SpaceQuotaResponse{}
}

func (mdMgr *SpaceQuotaMetadataManager) CreateResourceObject() common.IResource {
	return &SpaceQuotaResource{}
}

func (mdMgr *SpaceQuotaMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewSpaceQuotaMetadataById(guid)
}

func (mdMgr *SpaceQuotaMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*SpaceQuotaResponse)
	for _, item := range resp.Resources {
		item.Entity.Guid = item.Meta.Guid
		metadata := NewSpaceQuotaMetadata(item.Entity)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *SpaceQuotaMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*SpaceQuotaResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewSpaceQuotaMetadata(resourceType.Entity)
	return metadata
}
