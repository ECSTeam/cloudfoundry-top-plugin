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
)

type GlobalManager struct {
	appMdMgr *app.AppMetadataManager
	//orgMdMgr *OrgMetadataManager
	//spaceMdMgr *SpaceMetadataManager
	orgQuotaMdMgr   *orgQuota.OrgQuotaMetadataManager
	spaceQuotaMdMgr *spaceQuota.SpaceQuotaMetadataManager

	cliConnection plugin.CliConnection

	appDeleteQueue map[string]string
	refreshQueue   map[string]time.Time
	refreshNow     chan bool
	refreshLock    sync.Mutex

	refreshAppInstanceStatisticsQueue map[string]time.Time
	refreshAppInstanceStatisticsNow   chan bool
	refreshAppInstanceStatisticsLock  sync.Mutex

	loadMetadataInProgress bool
}

func NewGlobalManager(conn plugin.CliConnection) *GlobalManager {

	mgr := &GlobalManager{}

	mgr.appMdMgr = app.NewAppMetadataManager()
	mgr.orgQuotaMdMgr = orgQuota.NewOrgQuotaMetadataManager(mgr)
	mgr.spaceQuotaMdMgr = spaceQuota.NewSpaceQuotaMetadataManager(mgr)

	mgr.cliConnection = conn

	mgr.appDeleteQueue = make(map[string]string)
	mgr.refreshQueue = make(map[string]time.Time)
	mgr.refreshNow = make(chan bool, 2)

	mgr.refreshAppInstanceStatisticsQueue = make(map[string]time.Time)
	mgr.refreshAppInstanceStatisticsNow = make(chan bool, 2)

	// Set set the time of event data end date/time here so we don't end up loading
	// events after we've already started counting them from the firehose.
	now := time.Now()
	crashData.LoadEventsUntilTime = &now

	go mgr.loadMetadataThread()
	go mgr.loadMetadataAppStatisticsThread()

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

func (mgr *GlobalManager) RequestRefreshAppInstanceStatisticsMetadata(appId string) {
	mgr.refreshAppInstanceStatisticsLock.Lock()
	mgr.refreshAppInstanceStatisticsQueue[appId] = time.Now()
	mgr.refreshAppInstanceStatisticsLock.Unlock()
	mgr.wakeRefreshAppInstanceStatisticsThread()
}

func (mgr *GlobalManager) wakeRefreshAppInstanceStatisticsThread() {
	select {
	case mgr.refreshAppInstanceStatisticsNow <- true:
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
			appMetadata := mgr.appMdMgr.FindAppMetadataInternal(appId, false)
			timeSinceLastLoad := time.Now().Sub(appMetadata.CacheTime)
			appName := appMetadata.Name
			toplog.Debug("Metadata - appId: %v name: [%v] - inqueue check time since last load: %v", appId, appName, timeSinceLastLoad)
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug("Metadata - appId: %v name: [%v] - Needs to be loaded now", appId, appName)
				newAppMetadata, err := mgr.appMdMgr.GetAppMetadataInternal(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn("Metadata - appId: %v name: [%v] - Error: %v", appId, appName, err)
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
					}
					mgr.refreshLock.Unlock()

				}
			} else {
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

func (mgr *GlobalManager) loadMetadataAppStatisticsThread() {

	minimumLoadTimeMS := time.Millisecond * 1000
	veryLongtime := time.Hour * 10000
	minNextLoadTime := veryLongtime

	for {

		toplog.Debug("Metadata appInstanceStatistics - sleep time: %v", minNextLoadTime)

		select {
		case <-mgr.refreshAppInstanceStatisticsNow:
		case <-time.After(minNextLoadTime):
		}

		minNextLoadTime = veryLongtime
		toplog.Debug("Metadata appInstanceStatistics thread is awake")

		mgr.refreshAppInstanceStatisticsLock.Lock()
		queue := make([]string, 0)
		for appId, _ := range mgr.refreshAppInstanceStatisticsQueue {
			queue = append(queue, appId)
		}
		mgr.refreshAppInstanceStatisticsLock.Unlock()

		for _, appId := range queue {
			now := time.Now()
			appMetadata := mgr.appMdMgr.FindAppMetadataInternal(appId, false)
			appName := appMetadata.Name

			appInstanceStatistics := appStatistics.FindAppStatisticMetadataInternal(appId)
			timeSinceLastLoad := veryLongtime
			if appInstanceStatistics != nil {
				timeSinceLastLoad = time.Now().Sub(*appInstanceStatistics.CacheTime)
			}
			toplog.Debug("Metadata appInstanceStatistics - appId: %v name: [%v] - inqueue check time since last load: %v", appId, appName, timeSinceLastLoad)
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug("Metadata appInstanceStatistics - appId: %v name: [%v] - Needs to be loaded now", appId, appName)
				err := appStatistics.LoadAppStatisticCache(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn("Metadata appInstanceStatistics - appId: %v name: [%v] - Error: %v", appId, appName, err)
				} else {
					toplog.Info("Metadata appInstanceStatistics - appId: %v name: [%v] - Load complete", appId, appName)

					if mgr.anyAppInstanceHaveState(appId, "DOWN") {
						err := appInstances.LoadAppInstancesCache(mgr.cliConnection, appId)
						if err != nil {
							toplog.Warn("Metadata appInstances - appId: %v name: [%v] - Error: %v", appId, appName, err)
						} else {
							toplog.Info("Metadata appInstances - appId: %v name: [%v] - Load complete", appId, appName)
						}
					}

					// Only delete if request for reload was queue before we loaded the data
					// This prevents a timing issue a request is in progress while another one is queued.
					mgr.refreshAppInstanceStatisticsLock.Lock()
					queueTime := mgr.refreshAppInstanceStatisticsQueue[appId]
					if queueTime.Before(now) {
						toplog.Debug("Metadata appInstanceStatistics - appId: %v name: [%v] - Remove from queue queueTime: %v now: %v", appId, appName, queueTime, now)
						delete(mgr.refreshAppInstanceStatisticsQueue, appId)
					}
					mgr.refreshAppInstanceStatisticsLock.Unlock()

				}
			} else {
				toplog.Debug("Metadata appInstanceStatistics - appId %v name: [%v] - Too soon to reload", appId, appName)
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug("Metadata appInstanceStatistics - appId %v name: [%v] - Try to load in: %v", appId, appName, nextLoadTime)
				if minNextLoadTime > nextLoadTime {
					toplog.Debug("Metadata appInstanceStatistics - appId %v name: [%v] - value was min: %v", appId, appName, nextLoadTime)
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}

func (mgr *GlobalManager) anyAppInstanceHaveState(appId string, state string) bool {
	appInstanceStatistics := appStatistics.FindAppStatisticMetadataInternal(appId)
	if appInstanceStatistics == nil {
		return false
	}
	for _, stat := range appInstanceStatistics.Data {
		if stat.State == state {
			return true
		}
	}
	return false
}
