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
	"fmt"
	"sync"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/app"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/domain"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/org"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/route"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/space"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/stack"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"

	"code.cloudfoundry.org/cli/plugin"
)

type Manager struct {
	appMdMgr *app.AppMetadataManager
	//orgMdMgr *OrgMetadataManager
	//spaceMdMgr *SpaceMetadataManager

	mu sync.Mutex

	appDeleteQueue map[string]string

	refreshNow    chan bool
	refreshQueue  map[string]string
	cliConnection plugin.CliConnection
}

func NewManager(conn plugin.CliConnection) *Manager {

	mgr := &Manager{}

	mgr.appMdMgr = app.NewAppMetadataManager()

	mgr.appDeleteQueue = make(map[string]string)

	mgr.refreshQueue = make(map[string]string)
	mgr.refreshNow = make(chan bool)
	mgr.cliConnection = conn

	go mgr.loadMetadataThread()

	return mgr
}

func (mgr *Manager) GetAppMdManager() *app.AppMetadataManager {
	return mgr.appMdMgr
}

// Load all the metadata.  This is a blocking call.
func (mgr *Manager) LoadMetadata() {
	toplog.Info("Manager>loadMetadata")

	stack.LoadStackCache(mgr.cliConnection)

	mgr.appMdMgr.LoadAppCache(mgr.cliConnection)
	space.LoadSpaceCache(mgr.cliConnection)
	org.LoadOrgCache(mgr.cliConnection)

	route.LoadRouteCache(mgr.cliConnection)
	domain.LoadDomainCache(mgr.cliConnection)

}

func (mgr *Manager) IsAppDeleted(appId string) bool {
	return mgr.appDeleteQueue[appId] != ""
}

func (mgr *Manager) RemoveAppFromDeletedQueue(appId string) {
	delete(mgr.appDeleteQueue, appId)
}

// Request a refresh of specific app metadata
func (mgr *Manager) RequestRefreshAppMetadata(appId string) {
	mgr.refreshQueue[appId] = appId
	mgr.wakeRefreshThread()
}

func (mgr *Manager) wakeRefreshThread() {
	mgr.refreshNow <- true
}

func (mgr *Manager) loadMetadataThread() {

	minimumLoadTimeMS := time.Millisecond * 10000
	veryLongtime := time.Hour * 10000
	minNextLoadTime := veryLongtime

	for {

		toplog.Debug(fmt.Sprintf("Metadata - sleep time: %v", minNextLoadTime))

		select {
		case <-mgr.refreshNow:
			//mui.updateDisplay(g)
		case <-time.After(minNextLoadTime):
			//mui.updateDisplay(g)
		}

		minNextLoadTime = veryLongtime
		toplog.Debug("Metadata cache thread is awake")
		for _, appId := range mgr.refreshQueue {
			appMetadata := mgr.appMdMgr.FindAppMetadataInternal(appId, false)
			timeSinceLastLoad := time.Now().Sub(appMetadata.CacheTime)
			appName := appMetadata.Name
			toplog.Debug(fmt.Sprintf("Metadata - appId: %v name: [%v] - inqueue check time since last load: %v", appId, appName, timeSinceLastLoad))
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug(fmt.Sprintf("Metadata - appId: %v name: [%v] - Needs to be loaded now", appId, appName))
				newAppMetadata, err := mgr.appMdMgr.GetAppMetadataInternal(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn(fmt.Sprintf("Metadata - appId: %v name: [%v] - Error: %v", appId, appName, err))
				} else {
					toplog.Info(fmt.Sprintf("Metadata - appId: %v name: [%v] - Load start", appId, appName))
					if newAppMetadata.Name != "" {
						// Only save if it really loaded
						mgr.appMdMgr.GetAppMetadataMap()[appId] = newAppMetadata
					} else {
						// If we can't reload this appId the it must have been deleted
						// Remove from metadata cache AND remove from appstats in "current" processor
						delete(mgr.appMdMgr.GetAppMetadataMap(), appId)
						mgr.appDeleteQueue[appId] = appId
						toplog.Info(fmt.Sprintf("Metadata - appId: %v name: [%v] - Removed from cache as it doesn't seem to exist", appId, appName))
					}
					toplog.Info(fmt.Sprintf("Metadata - appId: %v name: [%v] - Load complete", appId, newAppMetadata.Name))
					delete(mgr.refreshQueue, appId)
				}
			} else {
				toplog.Debug(fmt.Sprintf("Metadata - appId %v name: [%v] - Too soon to reload", appId, appName))
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug(fmt.Sprintf("Metadata - appId %v name: [%v] - Try to load in: %v", appId, appName, nextLoadTime))
				if minNextLoadTime > nextLoadTime {
					toplog.Debug(fmt.Sprintf("Metadata - appId %v name: [%v] - value was min: %v", appId, appName, nextLoadTime))
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}
