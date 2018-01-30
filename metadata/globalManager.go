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
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
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
	appMdMgr           *app.AppMetadataManager
	appInstMdMgr       *appInstances.AppInstanceMetadataManager
	orgMdMgr           *org.OrgMetadataManager
	orgQuotaMdMgr      *orgQuota.OrgQuotaMetadataManager
	spaceMdMgr         *space.SpaceMetadataManager
	spaceQuotaMdMgr    *spaceQuota.SpaceQuotaMetadataManager
	stackMdMgr         *stack.StackMetadataManager
	isoSegMdMgr        *isolationSegment.IsolationSegmentMetadataManager
	domainSharedMdMgr  *domain.DomainSharedMetadataManager
	domainPrivateMdMgr *domain.DomainPrivateMetadataManager
	domainFinder       *domain.DomainFinder
	routeMdMgr         *route.RouteMetadataManager

	cliConnection plugin.CliConnection

	// Collection of appIds that are monitored for container changes
	// The time is when the app was last viewed -- it will be used for a TTL
	// If app detail hasn't been viewed for awhile, it will be removed from list
	monitoredAppDetails     map[string]*time.Time
	monitoredAppDetailsLock sync.Mutex

	statusMsg chan string

	loadMetadataInProgress bool

	loadHandler *common.LoadHandler
}

func NewGlobalManager(conn plugin.CliConnection, statusMsg chan string) *GlobalManager {

	mgr := &GlobalManager{statusMsg: statusMsg}

	mgr.loadHandler = common.NewLoadHandler(conn)

	mgr.appMdMgr = app.NewAppMetadataManager(mgr)
	mgr.appInstMdMgr = appInstances.NewAppInstanceMetadataManager(mgr)
	mgr.orgMdMgr = org.NewOrgMetadataManager(mgr)
	mgr.orgQuotaMdMgr = orgQuota.NewOrgQuotaMetadataManager(mgr)

	mgr.spaceMdMgr = space.NewSpaceMetadataManager(mgr)
	mgr.spaceQuotaMdMgr = spaceQuota.NewSpaceQuotaMetadataManager(mgr)

	mgr.stackMdMgr = stack.NewStackMetadataManager(mgr)
	mgr.isoSegMdMgr = isolationSegment.NewIsolationSegmentMetadataManager(mgr)

	mgr.domainSharedMdMgr = domain.NewDomainSharedMetadataManager(mgr)
	mgr.domainPrivateMdMgr = domain.NewDomainPrivateMetadataManager(mgr)
	mgr.domainFinder = domain.NewDomainFinder(mgr.domainSharedMdMgr, mgr.domainPrivateMdMgr)

	mgr.routeMdMgr = route.NewRouteMetadataManager(mgr)

	mgr.cliConnection = conn

	mgr.monitoredAppDetails = make(map[string]*time.Time)

	// Set set the time of event data end date/time here so we don't end up loading
	// events after we've already started counting them from the firehose.
	now := time.Now()
	crashData.LoadEventsUntilTime = &now

	return mgr
}

func (mgr *GlobalManager) SetStatus(status string) {
	mgr.statusMsg <- status
}

func (mgr *GlobalManager) GetAppMdManager() *app.AppMetadataManager {
	return mgr.appMdMgr
}

func (mgr *GlobalManager) GetAppMetadataFromUrl(url string) ([]common.IMetadata, error) {
	return mgr.appMdMgr.GetMetadataFromUrl(url)
}

func (mgr *GlobalManager) GetAppInstMdManager() *appInstances.AppInstanceMetadataManager {
	return mgr.appInstMdMgr
}

func (mgr *GlobalManager) GetOrgMdManager() *org.OrgMetadataManager {
	return mgr.orgMdMgr
}

func (mgr *GlobalManager) GetOrgQuotaMdManager() *orgQuota.OrgQuotaMetadataManager {
	return mgr.orgQuotaMdMgr
}

func (mgr *GlobalManager) GetSpaceMdManager() *space.SpaceMetadataManager {
	return mgr.spaceMdMgr
}

func (mgr *GlobalManager) GetStackMdManager() *stack.StackMetadataManager {
	return mgr.stackMdMgr
}

func (mgr *GlobalManager) GetSpaceQuotaMdManager() *spaceQuota.SpaceQuotaMetadataManager {
	return mgr.spaceQuotaMdMgr
}

func (mgr *GlobalManager) GetIsoSegMdManager() *isolationSegment.IsolationSegmentMetadataManager {
	return mgr.isoSegMdMgr
}

func (mgr *GlobalManager) GetDomainSharedMdManager() *domain.DomainSharedMetadataManager {
	return mgr.domainSharedMdMgr
}

func (mgr *GlobalManager) GetDomainPrivateMdManager() *domain.DomainPrivateMetadataManager {
	return mgr.domainPrivateMdMgr
}
func (mgr *GlobalManager) GetDomainFinder() *domain.DomainFinder {
	return mgr.domainFinder
}

func (mgr *GlobalManager) GetRouteMdManager() *route.RouteMetadataManager {
	return mgr.routeMdMgr
}

func (mgr *GlobalManager) GetCliConnection() plugin.CliConnection {
	return mgr.cliConnection
}

// Load all the metadata.  This is a blocking call.
func (mgr *GlobalManager) LoadMetadata() {
	toplog.Info("GlobalManager>loadMetadata")

	mgr.loadMetadataInProgress = true

	mgr.isoSegMdMgr.LoadAllItems()
	mgr.stackMdMgr.LoadAllItems()
	mgr.appMdMgr.LoadAllItems()

	//time.Sleep(time.Second * 60)

	mgr.spaceMdMgr.LoadAllItems()
	mgr.orgMdMgr.LoadAllItems()

	mgr.routeMdMgr.LoadAllItems()

	mgr.domainSharedMdMgr.LoadAllItems()
	mgr.domainPrivateMdMgr.LoadAllItems()
	crashData.LoadCrashDataCache(mgr.cliConnection)

	mgr.loadMetadataInProgress = false

}

func (mgr *GlobalManager) FlushCache() {
	appStatistics.Clear()
	mgr.appInstMdMgr.Clear()
	mgr.LoadMetadata()
	mgr.orgQuotaMdMgr.Clear()
	mgr.spaceQuotaMdMgr.Clear()
}

// Request a refresh of specific app metadata
func (mgr *GlobalManager) RequestLoadOfItem(dataType common.DataType, guid string) {
	mgr.loadHandler.RequestLoadOfItem(dataType, guid, 0*time.Second)
}

func (mgr *GlobalManager) RequestLoadOfAll(dataType common.DataType) {
	mgr.loadHandler.RequestLoadOfAll(dataType, 0*time.Second)
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
		//appInstances.ClearAppInstancesMetadata(appId)
		mgr.appInstMdMgr.DeleteItem(appId)
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
	mgr.loadHandler.RequestLoadOfItem(common.APP_INST, appId, 0*time.Second)
}
