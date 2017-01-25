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

package common

import (
	"encoding/json"
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/mdGlobalManagerInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type createResponseObject func() IResponse
type createMetadataEntityObject func(guid string) IMetadata
type processResponse func(IResponse, []IMetadata) []IMetadata

type MdCommonManager struct {
	metadataMap       map[string]IMetadata
	mdGlobalManager   mdGlobalManagerInterface.MdGlobalManagerInterface
	fullLoadCacheTime time.Time
	loadInProgress    bool
	mu                sync.Mutex

	url                        string
	createResponseObject       createResponseObject
	createMetadataEntityObject createMetadataEntityObject
	processResponse            processResponse
}

func NewMdCommonManager(
	mdGlobalManager mdGlobalManagerInterface.MdGlobalManagerInterface,
	url string,
	createResponseObject createResponseObject,
	createMetadataEntityObject createMetadataEntityObject,
	processResponse processResponse) *MdCommonManager {
	mdMgr := &MdCommonManager{
		mdGlobalManager:            mdGlobalManager,
		url:                        url,
		createResponseObject:       createResponseObject,
		createMetadataEntityObject: createMetadataEntityObject,
		processResponse:            processResponse}
	mdMgr.metadataMap = make(map[string]IMetadata)
	return mdMgr
}

func (mdMgr *MdCommonManager) All() []IMetadata {
	metadataArray := []IMetadata{}
	for _, metadata := range mdMgr.metadataMap {
		metadataArray = append(metadataArray, metadata)
	}
	return metadataArray
}

func (mdMgr *MdCommonManager) Find(guid string) IMetadata {
	metadata := mdMgr.metadataMap[guid]
	if metadata == nil {
		mdMgr.RequestLoadCacheIfOld()
		metadata = mdMgr.createMetadataEntityObject(guid)
		// We mark this metadata as 60 mins old so it will be refreshed with async load
		metadata.SetCacheTime(metadata.GetCacheTime().Add(-60 * time.Minute))
	}
	return metadata
}

// Flush cache.  This will not force a reload - must call LoadCache if immedently reload is desired
func (mdMgr *MdCommonManager) FlushCache() {
	mdMgr.metadataMap = make(map[string]IMetadata)
	mdMgr.fullLoadCacheTime = time.Time{}
}

func (mdMgr *MdCommonManager) RequestLoadCacheIfOld() {
	if time.Now().Sub(mdMgr.fullLoadCacheTime) > time.Hour*24 {
		mdMgr.LoadCacheAysnc()
	}
}

func (mdMgr *MdCommonManager) LoadCacheAysnc() {

	if mdMgr.loadInProgress {
		toplog.Info("MdCommonManager.LoadCacheAysnc %v loadInProgress", mdMgr.url)
		return
	}

	mdMgr.loadInProgress = true
	loadAsync := func() {
		toplog.Info("MdCommonManager.LoadCacheAysnc %v loadAsync thread started", mdMgr.url)
		mdMgr.LoadCache(mdMgr.mdGlobalManager.GetCliConnection())
		mdMgr.loadInProgress = false
		toplog.Info("MdCommonManager.LoadCacheAysnc %v loadAsync thread complete", mdMgr.url)
	}
	go loadAsync()
}

func (mdMgr *MdCommonManager) LoadCache(cliConnection plugin.CliConnection) {
	metadataArray, err := mdMgr.getMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** metadata error: %v", err.Error())
		return
	}

	metadataMap := make(map[string]IMetadata)
	for _, metadata := range metadataArray {
		//toplog.Info("From Map - %+v", metadata)
		metadataMap[metadata.GetGuid()] = metadata
	}

	mdMgr.mu.Lock()
	defer mdMgr.mu.Unlock()
	mdMgr.metadataMap = metadataMap
	mdMgr.fullLoadCacheTime = time.Now()
}

func (mdMgr *MdCommonManager) getMetadata(cliConnection plugin.CliConnection) ([]IMetadata, error) {
	return mdMgr.GetMetadataFromUrl(cliConnection, mdMgr.url)
}

func (mdMgr *MdCommonManager) GetMetadataFromUrl(cliConnection plugin.CliConnection, url string) ([]IMetadata, error) {

	metadataArray := []IMetadata{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		resp := mdMgr.createResponseObject()
		err := json.Unmarshal(outputBytes, &resp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadataArray, err
		}
		metadataArray = mdMgr.processResponse(resp, metadataArray)
		return resp, nil
	}

	err := CallPagableAPI(cliConnection, url, handleRequest)

	return metadataArray, err

}
