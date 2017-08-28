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

package orgQuota

import (
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/mdGlobalManagerInterface"
)

type OrgQuotaMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewOrgQuotaMetadataManager(mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface) *OrgQuotaMetadataManager {
	url := "/v2/quota_definitions"
	mdMgr := &OrgQuotaMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, url, mdMgr, true)

	return mdMgr
}

func (mdMgr *OrgQuotaMetadataManager) FindItem(appId string) *OrgQuotaMetadata {
	return mdMgr.FindItemInternal(appId, false, true).(*OrgQuotaMetadata)
}

func (mdMgr *OrgQuotaMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewOrgQuotaMetadataById(guid)
}

func (mdMgr *OrgQuotaMetadataManager) CreateResourceObject() common.IResource {
	return &OrgQuotaResource{}
}

func (mdMgr *OrgQuotaMetadataManager) CreateResponseObject() common.IResponse {
	return &OrgQuotaResponse{}
}

func (mdMgr *OrgQuotaMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewOrgQuotaMetadataById(guid)
}

func (mdMgr *OrgQuotaMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*OrgQuotaResponse)
	for _, item := range resp.Resources {
		item.Entity.Guid = item.Meta.Guid
		metadata := NewOrgQuotaMetadata(item.Entity)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *OrgQuotaMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(OrgQuotaResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewOrgQuotaMetadata(resourceType.Entity)
	return metadata
}
