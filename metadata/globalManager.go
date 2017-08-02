// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
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

package metadata

import (
	"sync"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/appInstances"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/appStatistics"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/crashData"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/isolationSegment"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/orgQuota"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/spaceQuota"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/config"
)

type GlobalManager struct {
	appMdMgr *app.AppMetadataManager
	//orgMdMgr *OrgMetadataManager
	//spaceMdMgr *SpaceMetadataManager
	orgQuotaMdMgr   *orgQuota.OrgQuotaMetadataManager
	spaceQuotaMdMgr *spaceQuota.SpaceQuotaMetadataManager

	cliConnection plugin.CliConnection

	// Collection of appIds that are monitored for container changes
	// The time is when the app was last viewed -- it will be used for a TTL
	// If app detail hasn't been viewed for awhile, it will be removed from list
	monitoredAppDetails map[string]*time.Time

	appDeleteQueue map[string]string
	refreshQueue   map[string]time.Time
	refreshNow     chan bool
	refreshLock    sync.Mutex

	refreshAppInstancesQueue map[string]time.Time
	refreshAppInstancesNow   chan bool
	refreshAppInstancesLock  sync.Mutex

	loadMetadataInProgress bool
}

func NewGlobalManager(conn plugin.CliConnection) *GlobalManager {

	mgr := &GlobalManager{}

	mgr.appMdMgr = app.NewAppMetadataManager()
	mgr.orgQuotaMdMgr = orgQuota.NewOrgQuotaMetadataManager(mgr)
	mgr.spaceQuotaMdMgr = spaceQuota.NewSpaceQuotaMetadataManager(mgr)

	mgr.cliConnection = conn

	mgr.monitoredAppDetails = make(map[string]*time.Time)

	mgr.appDeleteQueue = make(map[string]string)
	mgr.refreshQueue = make(map[string]time.Time)
	mgr.refreshNow = make(chan bool, 2)

	mgr.refreshAppInstancesQueue = make(map[string]time.Time)
	mgr.refreshAppInstancesNow = make(chan bool, 2)

	// Set set the time of event data end date/time here so we don't end up loading
	// events after we've already started counting them from the firehose.
	now := time.Now()
	crashData.LoadEventsUntilTime = &now

	go mgr.loadMetadataThread()
	go mgr.loadMetadataAppInstancesThread()

	return mgr
}

func (mgr *GlobalManager) GetAppMdManager() *app.AppMetadataManager {
	return mgr.appMdMgr
}

func (mgr *GlobalManager) GetOrgQuotaMdManager() *orgQuota.OrgQuotaMetadataManager {
	return mgr.orgQuotaMdMgr
}

func (mgr *GlobalManager) GetSpaceQuotaMdManager() *spaceQuota.SpaceQuotaMetadataManager {
	return mgr.spaceQuotaMdMgr
}

func (mgr *GlobalManager) GetCliConnection() plugin.CliConnection {
	return mgr.cliConnection
}

// Load all the metadata.  This is a blocking call.
func (mgr *GlobalManager) LoadMetadata() {
	toplog.Info("GlobalManager>loadMetadata")

	mgr.loadMetadataInProgress = true

	isolationSegment.LoadCache(mgr.cliConnection)
	stack.LoadStackCache(mgr.cliConnection)

	mgr.appMdMgr.LoadAppCache(mgr.cliConnection)

	//time.Sleep(time.Second * 60)

	space.LoadSpaceCache(mgr.cliConnection)
	org.LoadOrgCache(mgr.cliConnection)

	route.LoadRouteCache(mgr.cliConnection)
	domain.LoadDomainCache(mgr.cliConnection)
	crashData.LoadCrashDataCache(mgr.cliConnection)

	mgr.loadMetadataInProgress = false

}

func (mgr *GlobalManager) FlushCache() {
	appStatistics.Clear()
	appInstances.Clear()
	mgr.LoadMetadata()
	mgr.orgQuotaMdMgr.FlushCache()
	mgr.spaceQuotaMdMgr.FlushCache()
}

func (mgr *GlobalManager) IsAppDeleted(appId string) bool {
	return mgr.appDeleteQueue[appId] != ""
}

func (mgr *GlobalManager) RemoveAppFromDeletedQueue(appId string) {
	delete(mgr.appDeleteQueue, appId)
}

// Request a refresh of specific app metadata
func (mgr *GlobalManager) RequestRefreshAppMetadata(appId string) {
	mgr.refreshLock.Lock()
	mgr.refreshQueue[appId] = time.Now()
	mgr.refreshLock.Unlock()
	mgr.wakeRefreshThread()
}

func (mgr *GlobalManager) wakeRefreshThread() {
	select {
	case mgr.refreshNow <- true:
	default:
		//case <-time.After(1 * time.Nanosecond):
	}
}

// Indicate that we should actively monitor app details (container updates) for given appId
func (mgr *GlobalManager) MonitorAppDetails(appId string, lastViewed *time.Time) {
	mgr.refreshAppInstancesLock.Lock()
	defer mgr.refreshAppInstancesLock.Unlock()
	mgr.monitoredAppDetails[appId] = lastViewed
}

func (mgr *GlobalManager) IsMonitorAppDetails(appId string) bool {

	mgr.refreshAppInstancesLock.Lock()
	defer mgr.refreshAppInstancesLock.Unlock()

	lastViewed := mgr.monitoredAppDetails[appId]
	if lastViewed == nil {
		// AppId is not monitored
		return false
	}

	now := time.Now()
	lastViewedDuration := now.Sub(*lastViewed)
	if lastViewedDuration > (time.Second * config.MonitorAppDetailTTL) {
		// App Detail monitor TTL expired
		toplog.Info("Ignore refresh App Instance metadata - AppdId [%v] TTL expired", appId)
		appInstances.ClearAppInstancesMetadata(appId)
		delete(mgr.monitoredAppDetails, appId)
		return false
	}
	return true
}

func (mgr *GlobalManager) RequestRefreshAppInstancesMetadata(appId string) {

	if !mgr.IsMonitorAppDetails(appId) {
		toplog.Debug("Ignore refresh App Instance metadata - AppdId [%v] not monitored", appId)
		return
	}

	mgr.refreshAppInstancesLock.Lock()
	mgr.refreshAppInstancesQueue[appId] = time.Now()
	mgr.refreshAppInstancesLock.Unlock()

	mgr.wakeRefreshAppInstancesThread()
}

func (mgr *GlobalManager) wakeRefreshAppInstancesThread() {
	select {
	case mgr.refreshAppInstancesNow <- true:
	default:
		//case <-time.After(1 * time.Nanosecond):
	}
}

func (mgr *GlobalManager) loadMetadataThread() {

	minimumLoadTimeMS := time.Millisecond * 10000
	veryLongtime := time.Hour * 10000
	minNextLoadTime := veryLongtime

	for {

		toplog.Debug("Metadata - sleep time: %v", minNextLoadTime)

		select {
		case <-mgr.refreshNow:
			//mui.updateDisplay(g)
		case <-time.After(minNextLoadTime):
			//mui.updateDisplay(g)
		}

		minNextLoadTime = veryLongtime

		mgr.refreshLock.Lock()
		queue := make([]string, 0)
		for appId, _ := range mgr.refreshQueue {
			queue = append(queue, appId)
		}
		mgr.refreshLock.Unlock()

		toplog.Debug("Metadata cache thread is awake")
		for _, appId := range queue {
			now := time.Now()
			removedFromQueue := false
			appMetadata := mgr.appMdMgr.FindAppMetadataInternal(appId, false)
			timeSinceLastLoad := time.Now().Sub(appMetadata.CacheTime)
			appName := appMetadata.Name
			toplog.Debug("Metadata - appId: %v name: [%v] - inqueue check time since last load: %v", appId, appName, timeSinceLastLoad)
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug("Metadata - appId: %v name: [%v] - Needs to be loaded now", appId, appName)
				newAppMetadata, err := mgr.appMdMgr.GetAppMetadataInternal(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn("Metadata - appId: %v name: [%v] - Error: %v", appId, appName, err)
					// Since we had an error trying to load the metadata -- just remove from queue to prevent endless retries to load it
					mgr.refreshLock.Lock()
					delete(mgr.refreshQueue, appId)
					removedFromQueue = true
					mgr.refreshLock.Unlock()
				} else {
					toplog.Info("Metadata - appId: %v name: [%v] - Load start", appId, appName)
					if newAppMetadata.Name != "" {
						// Only save if it really loaded
						mgr.appMdMgr.AddAppMetadata(newAppMetadata)
					} else {
						// If we can't reload this appId the it must have been deleted
						// Remove from metadata cache AND remove from appstats in "current" processor
						mgr.appMdMgr.DeleteAppMetadata(appId)
						mgr.appDeleteQueue[appId] = appId
						toplog.Info("Metadata - appId: %v name: [%v] - Removed from cache as it doesn't seem to exist", appId, appName)
					}
					toplog.Info("Metadata - appId: %v name: [%v] - Load complete", appId, newAppMetadata.Name)

					// Only delete if request for reload was queue before we loaded the data
					// This prevents a timing issue a request is in progress while another one is queued.
					mgr.refreshLock.Lock()
					queueTime := mgr.refreshQueue[appId]
					if queueTime.Before(now) {
						toplog.Debug("Metadata - appId: %v name: [%v] - Remove from queue queueTime: %v now: %v", appId, appName, queueTime, now)
						delete(mgr.refreshQueue, appId)
						removedFromQueue = true
					}
					mgr.refreshLock.Unlock()

				}
			}

			if !removedFromQueue {
				toplog.Debug("Metadata - appId %v name: [%v] - Too soon to reload", appId, appName)
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug("Metadata - appId %v name: [%v] - Try to load in: %v", appId, appName, nextLoadTime)
				if minNextLoadTime > nextLoadTime {
					toplog.Debug("Metadata - appId %v name: [%v] - value was min: %v", appId, appName, nextLoadTime)
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}

func (mgr *GlobalManager) loadMetadataAppInstancesThread() {

	minimumLoadTimeMS := time.Millisecond * 1000
	veryLongtime := time.Hour * 10000
	minNextLoadTime := veryLongtime

	for {

		toplog.Debug("Metadata appInstances - sleep time: %v", minNextLoadTime)

		select {
		case <-mgr.refreshAppInstancesNow:
		case <-time.After(minNextLoadTime):
		}

		minNextLoadTime = veryLongtime
		toplog.Debug("Metadata appInstances thread is awake")

		mgr.refreshAppInstancesLock.Lock()
		queue := make([]string, 0)
		for appId, _ := range mgr.refreshAppInstancesQueue {
			queue = append(queue, appId)
		}
		mgr.refreshAppInstancesLock.Unlock()

		for _, appId := range queue {
			now := time.Now()
			removedFromQueue := false
			appMetadata := mgr.appMdMgr.FindAppMetadataInternal(appId, false)
			appName := appMetadata.Name

			appInsts := appInstances.FindAppInstancesMetadataInternal(appId)
			timeSinceLastLoad := veryLongtime
			if appInsts != nil {
				timeSinceLastLoad = time.Now().Sub(*appInsts.CacheTime)
			}
			toplog.Debug("Metadata appInstances - appId: %v name: [%v] - inqueue check time since last load: %v", appId, appName, timeSinceLastLoad)
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug("Metadata appInstances - appId: %v name: [%v] - Needs to be loaded now", appId, appName)
				err := appInstances.LoadAppInstancesCache(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn("Metadata appInstances - appId: %v name: [%v] - Error: %v", appId, appName, err)
					// Since we had an error trying to load the metadata -- just remove from queue to prevent endless retries to load it
					mgr.refreshAppInstancesLock.Lock()
					delete(mgr.refreshAppInstancesQueue, appId)
					removedFromQueue = true
					mgr.refreshAppInstancesLock.Unlock()
				} else {
					toplog.Info("Metadata appInstances - appId: %v name: [%v] - Load complete", appId, appName)

					// Only delete if request for reload was queue before we loaded the data
					// This prevents a timing issue a request is in progress while another one is queued.
					mgr.refreshAppInstancesLock.Lock()
					queueTime := mgr.refreshAppInstancesQueue[appId]

					if queueTime.Before(now) {
						toplog.Debug("Metadata appInstances - appId: %v name: [%v] - Remove from queue queueTime: %v now: %v", appId, appName, queueTime, now)
						delete(mgr.refreshAppInstancesQueue, appId)
						removedFromQueue = true
						// Check if we have a race condition where we were in the middle of reloading metadata when the TTL expired
						if mgr.monitoredAppDetails[appId] == nil {
							toplog.Info("Clear App Instance metadata - AppdId [%v] TTL expired", appId)
							appInstances.ClearAppInstancesMetadata(appId)
						}
					}
					mgr.refreshAppInstancesLock.Unlock()

				}
			}

			if !removedFromQueue {
				toplog.Debug("Metadata appInstances - appId %v name: [%v] - Too soon to reload", appId, appName)
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug("Metadata appInstances - appId %v name: [%v] - Try to load in: %v", appId, appName, nextLoadTime)
				if minNextLoadTime > nextLoadTime {
					toplog.Debug("Metadata appInstances - appId %v name: [%v] - value was min: %v", appId, appName, nextLoadTime)
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}

func (mgr *GlobalManager) anyAppInstanceHaveState(appId string, state string) bool {
	appInstances := appStatistics.FindAppStatisticMetadataInternal(appId)
	if appInstances == nil {
		return false
	}
	for _, stat := range appInstances.Data {
		if stat.State == state {
			return true
		}
	}
	return false
}
