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

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"

	"code.cloudfoundry.org/cli/plugin"
)

type Manager struct {
	appMdMgr *AppMetadataManager
	//orgMdMgr *OrgMetadataManager
	//spaceMdMgr *SpaceMetadataManager

	mu            sync.Mutex
	refreshNow    chan bool
	refreshQueue  map[string]string
	cliConnection plugin.CliConnection
}

func NewManager(conn plugin.CliConnection) *Manager {

	mgr := &Manager{}

	mgr.appMdMgr = NewAppMetadataManager()

	mgr.refreshQueue = make(map[string]string)
	mgr.refreshNow = make(chan bool)
	mgr.cliConnection = conn

	go mgr.loadMetadataThread()

	return mgr
}

func (mgr *Manager) GetAppMdManager() *AppMetadataManager {
	return mgr.appMdMgr
}

// Load all the metadata.  This is a blocking call.
func (mgr *Manager) LoadMetadata() {
	toplog.Info("Manager>loadMetadata")
	mgr.appMdMgr.LoadAppCache(mgr.cliConnection)
	LoadStackCache(mgr.cliConnection)
	LoadSpaceCache(mgr.cliConnection)
	LoadOrgCache(mgr.cliConnection)
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
			appMetadata := mgr.appMdMgr.appMetadataMap[appId]
			timeSinceLastLoad := time.Now().Sub(appMetadata.cacheTime)
			toplog.Debug(fmt.Sprintf("Metadata - appId: %v - inqueue check time since last load: %v", appId, timeSinceLastLoad))
			if timeSinceLastLoad > minimumLoadTimeMS {
				toplog.Debug(fmt.Sprintf("Metadata - appId: %v - Needs to be loaded now", appId))
				newAppMetadata, err := mgr.appMdMgr.getAppMetadata(mgr.cliConnection, appId)
				if err != nil {
					toplog.Warn(fmt.Sprintf("Metadata - appId: %v - Error: %v", appId, err))
				} else {
					toplog.Info(fmt.Sprintf("Metadata - appId: %v - Load start", appId))
					mgr.appMdMgr.appMetadataMap[appId] = newAppMetadata
					toplog.Info(fmt.Sprintf("Metadata - appId: %v - Load complete", appId))
				}
			} else {
				toplog.Debug(fmt.Sprintf("Metadata - appId %v - Too soon to reload", appId))
				nextLoadTime := minimumLoadTimeMS - timeSinceLastLoad
				toplog.Debug(fmt.Sprintf("Metadata - appId %v - Try to load in: %v", appId, nextLoadTime))
				if minNextLoadTime > nextLoadTime {
					toplog.Debug(fmt.Sprintf("Metadata - appId %v - value was min: %v", appId, nextLoadTime))
					minNextLoadTime = nextLoadTime
				}
			}
		}
	}
}
