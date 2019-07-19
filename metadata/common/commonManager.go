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

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

const DelayedRemovalFromCacheDuration = 15 * time.Second
const DefaultMinimumReloadDuration = time.Millisecond * 10000

type MetadataManager interface {
	NewItemById(guid string) IMetadata
	LoadItemInternal(guid string) (IMetadata, error)
}

type CommonMetadataManager struct {
	mdGlobalManager MdGlobalManagerInterface

	dataType DataType // APP, SPACE, ORG, etc
	url      string

	mm                 MetadataManager
	MetadataMapMutex   sync.Mutex
	MetadataMap        map[string]IMetadata
	autoLoadIfNotFound bool

	pendingDeleteFromCache map[string]*time.Time
	deletedFromCache       map[string]*time.Time

	minimumReloadDuration time.Duration
}

func NewCommonMetadataManager(
	mdGlobalManager MdGlobalManagerInterface,
	dataType DataType,
	url string,
	mm MetadataManager,
	minimumReloadDuration time.Duration) *CommonMetadataManager {
	commonMgr := &CommonMetadataManager{mdGlobalManager: mdGlobalManager, dataType: dataType, url: url, mm: mm}
	commonMgr.minimumReloadDuration = minimumReloadDuration
	commonMgr.clear()
	RegisterMetadataHandler(dataType, commonMgr)
	return commonMgr
}

func (commonMgr *CommonMetadataManager) GetUrl() string {
	return commonMgr.url
}

func (commonMgr *CommonMetadataManager) GetMdGlobalManager() MdGlobalManagerInterface {
	return commonMgr.mdGlobalManager
}

func (commonMgr *CommonMetadataManager) clear() {
	commonMgr.MetadataMap = make(map[string]IMetadata)
	commonMgr.pendingDeleteFromCache = make(map[string]*time.Time)
	commonMgr.deletedFromCache = make(map[string]*time.Time)
}

func (commonMgr *CommonMetadataManager) Clear() {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	commonMgr.clear()
}

func (commonMgr *CommonMetadataManager) CacheSize() int {
	return len(commonMgr.MetadataMap)
}

func (commonMgr *CommonMetadataManager) AddItem(metadataItem IMetadata) {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	commonMgr.MetadataMap[metadataItem.GetGuid()] = metadataItem
}

func (commonMgr *CommonMetadataManager) DeleteItem(guid string) {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	delete(commonMgr.MetadataMap, guid)
	delete(commonMgr.pendingDeleteFromCache, guid)
	now := time.Now()
	commonMgr.deletedFromCache[guid] = &now
}

func (commonMgr *CommonMetadataManager) FindItemInternal(guid string, requestLoadIfNotFound bool, createEmptyObjectIfNotFound bool) (IMetadata, bool) {

	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	found := true
	metadataItem := commonMgr.MetadataMap[guid]
	if metadataItem == nil {
		found = false
		if createEmptyObjectIfNotFound {
			metadataItem = commonMgr.mm.NewItemById(guid)
			now := time.Now()
			metadataItem.SetCacheTime(&now)
			if requestLoadIfNotFound {
				// TODO: Queue metadata load for this id
			} else {
				// We mark this metadata as 60 mins old
				//loadTime := appMetadata.CacheTime.Add(-60 * time.Minute)
				//appMetadata.CacheTime = &loadTime
			}
		}
	}
	return metadataItem, found
}

// Called via a seperate thread - after a delay, remove the requested guid from cache
func (commonMgr *CommonMetadataManager) DelayedRemovalFromCache(guid string, itemName string) {

	commonMgr.addToPendingDeleteFromCache(guid, itemName)
	time.Sleep(DelayedRemovalFromCacheDuration)
	toplog.Info("Metadata %v - guid: %v name: [%v] - Removed from cache as it doesn't seem to exist", commonMgr.dataType, guid, itemName)
	commonMgr.DeleteItem(guid)
}

func (commonMgr *CommonMetadataManager) addToPendingDeleteFromCache(guid string, itemName string) {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()

	if commonMgr.pendingDeleteFromCache[guid] != nil {
		// guid already queued for delete
		return
	}
	now := time.Now()
	commonMgr.pendingDeleteFromCache[guid] = &now
}

func (commonMgr *CommonMetadataManager) IsDeletedFromCache(guid string) bool {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	return commonMgr.deletedFromCache[guid] != nil
}

func (commonMgr *CommonMetadataManager) IsPendingDeleteFromCache(guid string) bool {
	commonMgr.MetadataMapMutex.Lock()
	defer commonMgr.MetadataMapMutex.Unlock()
	return commonMgr.pendingDeleteFromCache[guid] != nil
}

func (commonMgr *CommonMetadataManager) MetadataLoadMethod(guid string) error {
	return commonMgr.LoadItem(guid)
}

func (commonMgr *CommonMetadataManager) MinimumReloadDuration() time.Duration {
	return commonMgr.minimumReloadDuration
}

// Last time data was loaded or nil if never
func (commonMgr *CommonMetadataManager) LastLoadTime(dataKey string) *time.Time {
	item, _ := commonMgr.FindItemInternal(dataKey, false, false)
	if item != nil {
		return item.GetCacheTime()
	}
	return nil
}

func (commonMgr *CommonMetadataManager) LoadItem(guid string) error {

	metadataItem, _ := commonMgr.FindItemInternal(guid, commonMgr.autoLoadIfNotFound, true)
	itemName := metadataItem.GetName()

	if commonMgr.IsPendingDeleteFromCache(guid) {
		toplog.Info("Metadata - Ignore metadataItem Load request as its been queued for cache deletion. guid: %v name: [%v] - Load start", guid, itemName)
		return nil
	}

	toplog.Info("Metadata %v - guid: %v name: [%v] - Load start", commonMgr.dataType, guid, itemName)
	start := time.Now()
	newMetadata, err := commonMgr.mm.LoadItemInternal(guid)
	if err != nil {
		return err
	} else {
		itemName = newMetadata.GetName()
		if itemName != "" {
			// Only save if it really loaded
			commonMgr.AddItem(newMetadata)
		} else {
			// If we can't reload this guid then it must have been deleted
			// Remove from metadata cache AND remove from appstats in "current" processor
			go commonMgr.DelayedRemovalFromCache(guid, itemName)
			toplog.Info("Metadata %v - guid: %v name: [%v] - Queue remove from cache as it doesn't seem to exist", commonMgr.dataType, guid, itemName)
		}
		elapsed := time.Since(start)
		toplog.Info("Metadata %v - guid: %v name: [%v] - Load complete in %s", commonMgr.dataType, guid, itemName, elapsed)
	}
	return nil
}
