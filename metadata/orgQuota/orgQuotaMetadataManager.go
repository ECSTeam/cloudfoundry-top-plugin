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
	"encoding/json"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/mdGlobalManagerInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type OrgQuotaMetadataManager struct {
	metadataMap       map[string]*OrgQuotaMetadata
	mdGlobalManager   mdGlobalManagerInterface.MdGlobalManagerInterface
	fullLoadCacheTime time.Time
	loadInProgress    bool
	mu                sync.Mutex
}

func NewOrgQuotaMetadataManager(mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface) *OrgQuotaMetadataManager {
	mgr := &OrgQuotaMetadataManager{mdGlobalManager: mdGlobalManager}
	mgr.metadataMap = make(map[string]*OrgQuotaMetadata)
	return mgr
}

func (mdMgr *OrgQuotaMetadataManager) All() []*OrgQuotaMetadata {
	metadataArray := []*OrgQuotaMetadata{}
	for _, metadata := range mdMgr.metadataMap {
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *OrgQuotaMetadataManager) Find(orgQuotaGuid string) *OrgQuotaMetadata {
	metadata := mdMgr.metadataMap[orgQuotaGuid]
	if metadata == nil {
		mdMgr.RequestLoadCacheIfOld()
		metadata = NewOrgQuotaMetadataById(orgQuotaGuid)
		// We mark this metadata as 60 mins old so it will be refreshed with async load
		metadata.cacheTime = metadata.cacheTime.Add(-60 * time.Minute)
	}
	return metadata
}

func (mdMgr *OrgQuotaMetadataManager) RequestLoadCacheIfOld() {
	if time.Now().Sub(mdMgr.fullLoadCacheTime) > time.Hour*24 {
		mdMgr.LoadCacheAysnc()
	}
}

func (mdMgr *OrgQuotaMetadataManager) LoadCacheAysnc() {

	if mdMgr.loadInProgress {
		toplog.Info("OrgQuotaMetadataManager.LoadCacheAysnc loadInProgress")
		return
	}

	mdMgr.loadInProgress = true
	loadAsync := func() {
		toplog.Info("OrgQuotaMetadataManager.LoadCacheAysnc loadAsync thread started")
		mdMgr.LoadCache(mdMgr.mdGlobalManager.GetCliConnection())
		mdMgr.loadInProgress = false
		toplog.Info("OrgQuotaMetadataManager.LoadCacheAysnc loadAsync thread complete")
	}
	go loadAsync()
}

func (mdMgr *OrgQuotaMetadataManager) LoadCache(cliConnection plugin.CliConnection) {
	metadataArray, err := mdMgr.getMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** orgQuota metadata error: %v", err.Error())
		return
	}

	metadataMap := make(map[string]*OrgQuotaMetadata)
	for _, metadata := range metadataArray {
		toplog.Info("From Map - orgQuota id: %v name:%v  MemoryLimit: %v", metadata.Guid, metadata.Name, metadata.MemoryLimit)
		metadataMap[metadata.Guid] = metadata
	}

	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	mdMgr.metadataMap = metadataMap
	mdMgr.fullLoadCacheTime = time.Now()
}

func (mdMgr *OrgQuotaMetadataManager) getMetadata(cliConnection plugin.CliConnection) ([]*OrgQuotaMetadata, error) {
	return GetMetadataFromUrl(cliConnection, "/v2/quota_definitions")
}

func GetMetadataFromUrl(cliConnection plugin.CliConnection, url string) ([]*OrgQuotaMetadata, error) {

	metadataArray := []*OrgQuotaMetadata{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var resp OrgQuotaResponse
		err := json.Unmarshal(outputBytes, &resp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadataArray, err
		}
		for _, item := range resp.Resources {
			item.Entity.Guid = item.Meta.Guid
			OrgQuotaMetadata := NewOrgQuotaMetadata(item.Entity)
			metadataArray = append(metadataArray, OrgQuotaMetadata)
		}
		return resp, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadataArray, err

}
