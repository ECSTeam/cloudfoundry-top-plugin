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

package domain

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

type DomainSharedMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewDomainSharedMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *DomainSharedMetadataManager {
	url := "/v2/shared_domains"
	mdMgr := &DomainSharedMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, common.DOMAIN_SHARED, url, mdMgr, false)
	return mdMgr
}

func (mdMgr *DomainSharedMetadataManager) FindItem(guid string) *DomainMetadata {
	return mdMgr.FindItemInternal(guid, false, true).(*DomainMetadata)
}

func (mdMgr *DomainSharedMetadataManager) GetAll() []*DomainMetadata {
	mdMgr.MetadataMapMutex.Lock()
	defer mdMgr.MetadataMapMutex.Unlock()
	metadataArray := []*DomainMetadata{}
	for _, metadata := range mdMgr.MetadataMap {
		metadataArray = append(metadataArray, metadata.(*DomainMetadata))
	}
	return metadataArray
}

func (mdMgr *DomainSharedMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewDomainMetadataById(guid)
}

func (mdMgr *DomainSharedMetadataManager) CreateResponseObject() common.IResponse {
	return &DomainResponse{}
}

func (mdMgr *DomainSharedMetadataManager) CreateResourceObject() common.IResource {
	return &DomainResource{}
}

func (mdMgr *DomainSharedMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewDomainMetadataById(guid)
}

func (mdMgr *DomainSharedMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*DomainResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *DomainSharedMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*DomainResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewDomainMetadata(resourceType.Entity)
	metadata.SharedDomain = true
	return metadata
}
