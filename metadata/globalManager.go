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
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/loader"
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
	monitoredAppDetails     map[string]*time.Time
	monitoredAppDetailsLock sync.Mutex

	loadMetadataInProgress bool

	loadHandler *loader.LoadHandler
}

func NewGlobalManager(conn plugin.CliConnection) *GlobalManager {

	mgr := &GlobalManager{}

	mgr.loadHandler = loader.NewLoadHandler(conn)

	mgr.appMdMgr = app.NewAppMetadataManager()
	mgr.orgQuotaMdMgr = orgQuota.NewOrgQuotaMetadataManager(mgr)
	mgr.spaceQuotaMdMgr = spaceQuota.NewSpaceQuotaMetadataManager(mgr)

	mgr.cliConnection = conn

	mgr.monitoredAppDetails = make(map[string]*time.Time)

	// Set set the time of event data end date/time here so we don't end up loading
	// events after we've already started counting them from the firehose.
	now := time.Now()
	crashData.LoadEventsUntilTime = &now

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

	mgr.appMdMgr.LoadCache(mgr.cliConnection)

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

// Request a refresh of specific app metadata
func (mgr *GlobalManager) RequestRefreshAppMetadata(appId string) {
	mgr.loadHandler.RequestLoadOfItem(loader.APP, appId, 0*time.Second)
}

// Indicate that we should actively monitor app details (container updates) for given appId
func (mgr *GlobalManager) MonitorAppDetails(appId string, lastViewed *time.Time) {
	mgr.monitoredAppDetailsLock.Lock()
	defer mgr.monitoredAppDetailsLock.Unlock()
	mgr.monitoredAppDetails[appId] = lastViewed
}

func (mgr *GlobalManager) IsMonitorAppDetails(appId string) bool {

	mgr.monitoredAppDetailsLock.Lock()
	defer mgr.monitoredAppDetailsLock.Unlock()

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
	mgr.loadHandler.RequestLoadOfItem(loader.APP_INST, appId, 0*time.Second)
}
