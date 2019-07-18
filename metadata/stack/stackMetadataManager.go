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

import (
	"regexp"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type StackMetadataManager struct {
	*common.CommonV2ResponseManager
	stackGroups []*StackGroup
}

func (mdMgr *StackMetadataManager) FindStackGroup(stackGroupId string) *StackGroup {
	for _, stackGroup := range mdMgr.stackGroups {
		if stackGroup.Guid == stackGroupId {
			return stackGroup
		}
	}
	return nil
}

func (mdMgr *StackMetadataManager) FindStackGroupByStackGuid(stackGuid string) *StackGroup {
	for _, stackGroup := range mdMgr.stackGroups {
		for _, stackId := range stackGroup.StackIds {
			if stackId == stackGuid {
				return stackGroup
			}
		}
	}
	return nil
}

func (mdMgr *StackMetadataManager) GetAllStackGroups() []*StackGroup {
	return mdMgr.stackGroups
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
	mdMgr.MetadataMapMutex.Lock()
	defer mdMgr.MetadataMapMutex.Unlock()
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

func (mdMgr *StackMetadataManager) PostProcessLoad(metadataArray []common.IMetadata, err error) {

	// TODO: SG
	// TODO: Do we create a new StackGroup object if we don't find a match for a given Stack?
	// If so, maybe we need to invert the for loops so we find a match for each item in metadataArray
	// and if there isn't a StackGroup match, we create a new StackGroup object for it.

	stackGroups := mdMgr.createKnownStackGroups()

	for _, metadata := range metadataArray {
		found := false
		for _, stackGroup := range stackGroups {
			r := regexp.MustCompile(stackGroup.MatchNames)

			//stackMetadata := metadata.(*StackMetadata)
			//toplog.Info("metadata: %v  stackGroup: %v", metadata.GetName(), stackGroup)

			index := r.FindStringSubmatchIndex(metadata.GetName())
			if len(index) > 0 {
				stackGroup.StackIds = append(stackGroup.StackIds, metadata.GetGuid())
				found = true
				toplog.Debug("stack found: %v  stackGroup: %+v", metadata.GetName(), stackGroup)
				break
			}
		}
		if !found {
			stackGroup := &StackGroup{Guid: metadata.GetGuid(), Name: metadata.GetName(),
				MatchNames: metadata.GetName(), StackIds: []string{metadata.GetGuid()}}
			stackGroups = append(stackGroups, stackGroup)
			toplog.Debug("stack new: %v  stackGroup: %+v", metadata.GetName(), stackGroup)
		}
	}

	mdMgr.stackGroups = stackGroups

}

func (mdMgr *StackMetadataManager) createKnownStackGroups() []*StackGroup {

	stackGroups := make([]*StackGroup, 0)

	stackGroup := &StackGroup{Guid: "1", Name: "cflinuxfs", MatchNames: "cflinuxfs.*", StackIds: make([]string, 0)}
	stackGroups = append(stackGroups, stackGroup)

	//stackGroup := &StackGroup{Guid: "2", Name: "windows", MatchNames: "windows2019|windows", StackIds: make([]string)}
	//stackGroups = append(stackGroups, stackGroup)

	return stackGroups
}
