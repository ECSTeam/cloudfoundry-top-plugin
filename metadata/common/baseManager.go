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
	"sync"
	"time"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

const DelayedRemovalFromCacheDuration = 15 * time.Second

type MetadataManager interface {
	NewItemById(guid string) BaseMetadataItemI
	//LoadItem(cliConnection plugin.CliConnection, appId string) error
	LoadItemInternal(cliConnection plugin.CliConnection, appId string) (BaseMetadataItemI, error)
	LoadInternal(cliConnection plugin.CliConnection) ([]BaseMetadataItemI, error)
}

type BaseMetadataManager struct {
	MetadataManager
	mu          sync.Mutex
	MetadataMap map[string]BaseMetadataItemI

	pendingDeleteFromCache map[string]*time.Time
	deletedFromCache       map[string]*time.Time
}

func NewBaseMetadataManager(mdMgr MetadataManager) *BaseMetadataManager {
	commonMgr := &BaseMetadataManager{MetadataManager: mdMgr}
	commonMgr.MetadataMap = make(map[string]BaseMetadataItemI)

	commonMgr.pendingDeleteFromCache = make(map[string]*time.Time)
	commonMgr.deletedFromCache = make(map[string]*time.Time)

	return commonMgr
}

func (commonMgr *BaseMetadataManager) CacheSize() int {
	return len(commonMgr.MetadataMap)
}

func (commonMgr *BaseMetadataManager) AddItem(metadataItem BaseMetadataItemI) {
	commonMgr.mu.Lock()
	defer commonMgr.mu.Unlock()
	commonMgr.MetadataMap[metadataItem.GetGuid()] = metadataItem
}

func (commonMgr *BaseMetadataManager) DeleteItem(guid string) {
	commonMgr.mu.Lock()
	defer commonMgr.mu.Unlock()
	delete(commonMgr.MetadataMap, guid)
	delete(commonMgr.pendingDeleteFromCache, guid)
	now := time.Now()
	commonMgr.deletedFromCache[guid] = &now
}

func (commonMgr *BaseMetadataManager) FindItemInternal(guid string, requestLoadIfNotFound bool) BaseMetadataItemI {

	commonMgr.mu.Lock()
	defer commonMgr.mu.Unlock()

	//TODO: error: concurrent map read and map write
	metadataItem := commonMgr.MetadataMap[guid]
	if metadataItem == nil {
		metadataItem = commonMgr.NewItemById(guid)
		if requestLoadIfNotFound {
			// TODO: Queue metadata load for this id
		} else {
			// We mark this metadata as 60 mins old
			//loadTime := appMetadata.CacheTime.Add(-60 * time.Minute)
			//appMetadata.CacheTime = &loadTime
		}
	}
	return metadataItem
}

// Called via a seperate thread - after a delay, remove the requested guid from cache
func (commonMgr *BaseMetadataManager) DelayedRemovalFromCache(guid string, itemName string) {

	commonMgr.addToPendingDeleteFromCache(guid, itemName)
	time.Sleep(DelayedRemovalFromCacheDuration)
	toplog.Info("Metadata - guid: %v name: [%v] - Removed from cache as it doesn't seem to exist", guid, itemName)
	commonMgr.DeleteItem(guid)
}

func (commonMgr *BaseMetadataManager) addToPendingDeleteFromCache(guid string, itemName string) {
	commonMgr.mu.Lock()
	defer commonMgr.mu.Unlock()

	if commonMgr.pendingDeleteFromCache[guid] != nil {
		// guid already queued for delete
		return
	}
	now := time.Now()
	commonMgr.pendingDeleteFromCache[guid] = &now
}

func (commonMgr *BaseMetadataManager) IsDeletedFromCache(guid string) bool {
	return commonMgr.deletedFromCache[guid] != nil
}

func (commonMgr *BaseMetadataManager) IsPendingDeleteFromCache(guid string) bool {
	return commonMgr.pendingDeleteFromCache[guid] != nil
}

func (commonMgr *BaseMetadataManager) MetadataLoadMethod(cliConnection plugin.CliConnection, guid string) error {
	return commonMgr.LoadItem(cliConnection, guid)
}

func (commonMgr *BaseMetadataManager) MinimumReloadDuration() time.Duration {
	return time.Millisecond * 10000
}

// Last time data was loaded or nil if never
func (commonMgr *BaseMetadataManager) LastLoadTime(dataKey string) *time.Time {
	item := commonMgr.FindItemInternal(dataKey, false)
	if item != nil {
		return item.GetCacheTime()
	}
	return nil
}

func (commonMgr *BaseMetadataManager) LoadCache(cliConnection plugin.CliConnection) {
	metadataItemArray, err := commonMgr.LoadInternal(cliConnection)
	if err != nil {
		toplog.Warn("*** app metadata error: %v", err.Error())
		return
	}

	metadataMap := make(map[string]BaseMetadataItemI)
	for _, metadataItem := range metadataItemArray {
		//toplog.Debug("From Map - app id: %v name:%v", appMetadata.Guid, appMetadata.Name)
		metadataMap[metadataItem.GetGuid()] = metadataItem
	}

	commonMgr.MetadataMap = metadataMap
}

func (commonMgr *BaseMetadataManager) LoadItem(cliConnection plugin.CliConnection, guid string) error {

	metadataItem := commonMgr.FindItemInternal(guid, false)
	itemName := metadataItem.GetName()

	if commonMgr.IsPendingDeleteFromCache(guid) {
		toplog.Info("Metadata - Ignore metadataItem Load request as its been queued for cache deletion. guid: %v name: [%v] - Load start", guid, itemName)
		return nil
	}

	toplog.Info("Metadata - guid: %v name: [%v] - Load start", guid, itemName)
	newAppMetadata, err := commonMgr.LoadItemInternal(cliConnection, guid)
	if err != nil {
		return err
	} else {
		itemName = newAppMetadata.GetName()
		if itemName != "" {
			// Only save if it really loaded
			commonMgr.AddItem(newAppMetadata)
		} else {
			// If we can't reload this guid then it must have been deleted
			// Remove from metadata cache AND remove from appstats in "current" processor
			go commonMgr.DelayedRemovalFromCache(guid, itemName)
			toplog.Info("Metadata - guid: %v name: [%v] - Queue remove from cache as it doesn't seem to exist", guid, itemName)
		}
		toplog.Info("Metadata - guid: %v name: [%v] - Load complete", guid, itemName)
	}
	return nil
}
