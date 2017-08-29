// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.Domain/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type DomainPrivateMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewDomainPrivateMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *DomainPrivateMetadataManager {
	url := "/v2/private_domains"
	mdMgr := &DomainPrivateMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, url, mdMgr, true)

	return mdMgr
}

func (mdMgr *DomainPrivateMetadataManager) FindItem(appId string) *DomainMetadata {
	return mdMgr.FindItemInternal(appId, false, true).(*DomainMetadata)
}

func (mdMgr *DomainPrivateMetadataManager) GetAll() []*DomainMetadata {
	// TODO: Need to use parent lock
	//mdMgr.mu.Lock()
	//defer mdMgr.mu.Unlock()
	metadataArray := []*DomainMetadata{}
	for _, metadata := range mdMgr.MetadataMap {
		metadataArray = append(metadataArray, metadata.(*DomainMetadata))
	}
	return metadataArray
}

func (mdMgr *DomainPrivateMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewDomainMetadataById(guid)
}

func (mdMgr *DomainPrivateMetadataManager) CreateResponseObject() common.IResponse {
	return &DomainResponse{}
}

func (mdMgr *DomainPrivateMetadataManager) CreateResourceObject() common.IResource {
	return &DomainResource{}
}

func (mdMgr *DomainPrivateMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewDomainMetadataById(guid)
}

func (mdMgr *DomainPrivateMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*DomainResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *DomainPrivateMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*DomainResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewDomainMetadata(resourceType.Entity)
	return metadata
}

func (mdMgr *DomainPrivateMetadataManager) AddDomainMetadata(domainName string) *DomainMetadata {
	domain := NewDomainMetadataById(util.Pseudo_uuid())
	domain.Name = domainName
	mdMgr.AddItem(domain)
	return domain
}
