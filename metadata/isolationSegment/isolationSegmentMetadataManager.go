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

package isolationSegment

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

type IsolationSegmentMetadataManager struct {
	*common.CommonV2ResponseManager
}

func NewIsolationSegmentMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *IsolationSegmentMetadataManager {
	url := "/v3/isolation_segments"
	mdMgr := &IsolationSegmentMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, common.ISO_SEG, url, mdMgr, false)
	return mdMgr
}

func (mdMgr *IsolationSegmentMetadataManager) FindItem(guid string) *IsolationSegmentMetadata {
	if guid == "" {
		return UnknownIsolationSegment
	}
	if guid == DefaultIsolationSegmentGuid {
		return SharedIsolationSegment
	}
	return mdMgr.FindItemInternal(guid, false, true).(*IsolationSegmentMetadata)
}

func (mdMgr *IsolationSegmentMetadataManager) GetAll() []*IsolationSegmentMetadata {
	// TODO: Need to use parent lock
	//mdMgr.mu.Lock()
	//defer mdMgr.mu.Unlock()
	metadataArray := []*IsolationSegmentMetadata{}
	for _, metadata := range mdMgr.MetadataMap {
		metadataArray = append(metadataArray, metadata.(*IsolationSegmentMetadata))
	}
	return metadataArray
}

func (mdMgr *IsolationSegmentMetadataManager) LoadAllItems() {
	mdMgr.CommonV2ResponseManager.LoadAllItems()
	SharedIsolationSegment = mdMgr.findMetadataByName(SharedIsolationSegmentName)
	if SharedIsolationSegment == nil {
		// If we didn't find the shared segment, this must be a pre-isolation segment version of cloud foundry
		SharedIsolationSegment = NewIsolationSegmentMetadata(IsolationSegment{EntityCommon: common.EntityCommon{Guid: DefaultIsolationSegmentGuid}, Name: SharedIsolationSegmentName})
	}
}

func (mdMgr *IsolationSegmentMetadataManager) findMetadataByName(name string) *IsolationSegmentMetadata {
	if name == "" {
		return nil
	}
	for _, isoSeg := range mdMgr.MetadataMap {
		if isoSeg.GetName() == name {
			return isoSeg.(*IsolationSegmentMetadata)
		}
	}
	return nil
}

func (mdMgr *IsolationSegmentMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewIsolationSegmentMetadataById(guid)
}

func (mdMgr *IsolationSegmentMetadataManager) CreateResponseObject() common.IResponse {
	return &IsolationSegmentResponse{}
}

func (mdMgr *IsolationSegmentMetadataManager) CreateResourceObject() common.IResource {
	return &IsolationSegmentResponse{}
}

func (mdMgr *IsolationSegmentMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewIsolationSegmentMetadataById(guid)
}

func (mdMgr *IsolationSegmentMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*IsolationSegmentResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *IsolationSegmentMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*IsolationSegment)
	metadata := NewIsolationSegmentMetadata(*resourceType)
	return metadata
}

func (mdMgr *IsolationSegmentMetadataManager) GetNextUrl(response common.IResponse) string {
	isoSegresponse := response.(*IsolationSegmentResponse)
	nextUrl := isoSegresponse.Pagination.Next.Href
	return nextUrl
}
