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

type OrgQuotaMetadataManagerNew struct {
	*common.MdCommonManager
}

func NewOrgQuotaMetadataManager(mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface) *OrgQuotaMetadataManagerNew {
	url := "/v2/quota_definitions"
	mdMgr := &OrgQuotaMetadataManagerNew{}
	commonMgr := common.NewMdCommonManager(mdGlobalManager, url, mdMgr.createResponseObject, mdMgr.createMetadataEntityObject, mdMgr.processResponse)
	mdMgr.MdCommonManager = commonMgr
	return mdMgr
}

func (mdMgr *OrgQuotaMetadataManagerNew) createResponseObject() common.IResponse {
	return &OrgQuotaResponse{}
}

func (mdMgr *OrgQuotaMetadataManagerNew) createMetadataEntityObject(guid string) common.IMetadata {
	return NewOrgQuotaMetadataById(guid)
}

func (mdMgr *OrgQuotaMetadataManagerNew) processResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*OrgQuotaResponse)
	for _, item := range resp.Resources {
		item.Entity.Guid = item.Meta.Guid
		metadata := NewOrgQuotaMetadata(item.Entity)
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *OrgQuotaMetadataManagerNew) Find(guid string) *OrgQuotaMetadata {
	return mdMgr.MdCommonManager.Find(guid).(*OrgQuotaMetadata)
}
