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

package stack

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

type StackMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewStackMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *StackMetadataManager {
	url := "/v2/stacks"
	mdMgr := &StackMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, common.STACK, url, mdMgr, false)
	return mdMgr
}

func (mdMgr *StackMetadataManager) FindItem(guid string) *StackMetadata {
	return mdMgr.FindItemInternal(guid, false, true).(*StackMetadata)
}

func (mdMgr *StackMetadataManager) GetAll() []*StackMetadata {
	// TODO: Need to use parent lock
	//mdMgr.mu.Lock()
	//defer mdMgr.mu.Unlock()
	metadataArray := []*StackMetadata{}
	for _, metadata := range mdMgr.MetadataMap {
		metadataArray = append(metadataArray, metadata.(*StackMetadata))
	}
	return metadataArray
}

func (mdMgr *StackMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewStackMetadataById(guid)
}

func (mdMgr *StackMetadataManager) CreateResponseObject() common.IResponse {
	return &StackResponse{}
}

func (mdMgr *StackMetadataManager) CreateResourceObject() common.IResource {
	return &StackResource{}
}

func (mdMgr *StackMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewStackMetadataById(guid)
}

func (mdMgr *StackMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*StackResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *StackMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*StackResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewStackMetadata(resourceType.Entity)
	return metadata
}
