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

package org

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

type OrgMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewOrgMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *OrgMetadataManager {
	url := "/v2/organizations"
	mdMgr := &OrgMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, common.ORG, url, mdMgr, true)
	return mdMgr
}

func (mdMgr *OrgMetadataManager) FindItem(guid string) *OrgMetadata {
	return mdMgr.FindItemInternal(guid, false, true).(*OrgMetadata)
}

func (mdMgr *OrgMetadataManager) GetAll() []*OrgMetadata {
	// TODO: Need to use parent lock
	//mdMgr.mu.Lock()
	//defer mdMgr.mu.Unlock()
	metadataArray := []*OrgMetadata{}
	for _, metadata := range mdMgr.MetadataMap {
		metadataArray = append(metadataArray, metadata.(*OrgMetadata))
	}
	return metadataArray
}

func (mdMgr *OrgMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewOrgMetadataById(guid)
}

func (mdMgr *OrgMetadataManager) CreateResponseObject() common.IResponse {
	return &OrgResponse{}
}

func (mdMgr *OrgMetadataManager) CreateResourceObject() common.IResource {
	return &OrgResource{}
}

func (mdMgr *OrgMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewOrgMetadataById(guid)
}

func (mdMgr *OrgMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*OrgResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *OrgMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*OrgResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewOrgMetadata(resourceType.Entity)
	return metadata
}
