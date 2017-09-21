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

	"code.cloudfoundry.org/cli/plugin"
)

/*
TODO:

Need to deal with minimumLoadTime -- timeSinceLastLoad > minimumLoadTimeMS
  We don't have a generic way to get the last load time for items / all of a dataType


*/

const Infinity time.Duration = 1<<63 - 1
const MaxLoadAttempts = 5
const WaitToReloadOnErrorDuration = time.Second * 30
const ALL = "ALL"

var metadataHandlerMap map[DataType]MetadataHandler = make(map[DataType]MetadataHandler)

type MetadataHandler interface {
	MetadataLoadMethod(guid string) error
	MinimumReloadDuration() time.Duration
	LastLoadTime(dataKey string) *time.Time
}

type LoadHandler struct {
	cliConnection plugin.CliConnection

	//loadRequestQueueX map[DataType]map[string]*LoadRequest

	loadRequestQueue []*LoadRequest

	loadWake chan bool
	loadLock sync.Mutex
}

type LoadRequest struct {
	dataType     DataType   // APP, SPACE, ORG, etc
	dataKey      string     // GUID of data (e.g., app/space/org)
	loadAfter    *time.Time // Load data after this time (or now if nil)
	loadAttempts int        // Number of times we've attempted to load this data
}

func RegisterMetadataHandler(dataType DataType, metadataHandler MetadataHandler) {
	toplog.Info("Registering metadata load handler for: %v", dataType)
	metadataHandlerMap[dataType] = metadataHandler
}

func NewLoadHandler(conn plugin.CliConnection) *LoadHandler {

	loadHandler := &LoadHandler{}
	loadHandler.cliConnection = conn

	loadHandler.loadRequestQueue = make([]*LoadRequest, 0)
	loadHandler.loadWake = make(chan bool, 2)

	go loadHandler.loadThread()

	return loadHandler
}

func (lh *LoadHandler) RequestLoadOfItem(dataType DataType, guid string, delayBeforeLoad time.Duration) {
	lh.requestLoadInternal(dataType, guid, delayBeforeLoad)
}

func (lh *LoadHandler) RequestLoadOfAll(dataType DataType, delayBeforeLoad time.Duration) {
	lh.requestLoadInternal(dataType, ALL, delayBeforeLoad)
}

func (lh *LoadHandler) isAlreadyQueued(dataType DataType, guid string) bool {
	for _, loadRequest := range lh.loadRequestQueue {
		if loadRequest.dataType == dataType && (loadRequest.dataKey == guid || loadRequest.dataKey == ALL) {
			return true
		}
	}
	return false
}

func (lh *LoadHandler) removeIndividualRequests(dataType DataType) {
	queue := lh.loadRequestQueue
	for i := len(queue) - 1; i >= 0; i-- {
		item := queue[i]
		if item.dataType == dataType && item.dataKey != ALL {
			queue = append(queue[:i], queue[i+1:]...)
		}
	}
	lh.loadRequestQueue = queue
}

func (lh *LoadHandler) requestLoadInternal(dataType DataType, guid string, delayBeforeLoad time.Duration) {
	now := time.Now()
	loadAfter := now.Add(delayBeforeLoad)
	loadRequest := &LoadRequest{dataType: dataType, dataKey: guid, loadAfter: &loadAfter}
	lh.RequestLoad(loadRequest)
}

func (lh *LoadHandler) RequestLoad(loadRequest *LoadRequest) {
	lh.loadLock.Lock()
	defer lh.loadLock.Unlock()

	if lh.isAlreadyQueued(loadRequest.dataType, loadRequest.dataKey) {
		return
	}

	if loadRequest.dataKey == ALL {
		lh.removeIndividualRequests(loadRequest.dataType)
	}

	toplog.Debug("Metadata loader - RequestLoad queued. loadRequest [%+v]", loadRequest)

	lh.loadRequestQueue = append(lh.loadRequestQueue, loadRequest)

	toplog.Debug("Metadata loader - queue len [%v]", len(lh.loadRequestQueue))

	lh.wakeLoadThread()
}

func (lh *LoadHandler) adjustLoadTimeIfNeeded(loadRequest *LoadRequest) bool {
	metadataHandler := metadataHandlerMap[loadRequest.dataType]
	if metadataHandler == nil {
		toplog.Error("Metadata loader - MetadataLoadHandler not found: %v", loadRequest)
		return false
	}
	lastLoadTime := metadataHandler.LastLoadTime(loadRequest.dataKey)
	if lastLoadTime != nil {
		now := time.Now()
		lastReloadDuration := now.Sub(*lastLoadTime)
		minReloadDuration := metadataHandler.MinimumReloadDuration()
		if lastReloadDuration < minReloadDuration {

			toplog.Info("Metadata loader - Adjusting load time. lastReloadDuration: %v minReloadDuration: %v", lastReloadDuration, minReloadDuration)

			requestedReloadDuration := loadRequest.loadAfter.Sub(now)
			if requestedReloadDuration < minReloadDuration {
				adjustedNextLoadTime := lastLoadTime.Add(minReloadDuration)
				loadRequest.loadAfter = &adjustedNextLoadTime
				return true
			}
		}
	}
	return false
}

func (lh *LoadHandler) wakeLoadThread() {
	select {
	case lh.loadWake <- true:
	default:
	}
}

func (lh *LoadHandler) loadThread() {

	for {

		minNextLoadTime := lh.findMinimumNextLoadDuration()
		toplog.Debug("Metadata loader - sleep time: %v", minNextLoadTime)

		select {
		case <-lh.loadWake:
		case <-time.After(minNextLoadTime):
		}

		readyToLoad := lh.getRequestReadyToLoad()
		toplog.Debug("Metadata loader - Wokeup to load: %+v", readyToLoad)

		lh.load(readyToLoad)
	}
}

func (lh *LoadHandler) getRequestReadyToLoad() *LoadRequest {

	lh.loadLock.Lock()
	defer lh.loadLock.Unlock()
	now := time.Now()

	// Loop through queue and find an item that is ready to load
	// delete item from queue (array)
	var readyToLoad *LoadRequest
	queue := lh.loadRequestQueue
	for i := len(queue) - 1; i >= 0; i-- {
		item := queue[i]
		if now.After(*item.loadAfter) {
			if !lh.adjustLoadTimeIfNeeded(item) {
				readyToLoad = item
				queue = append(queue[:i], queue[i+1:]...)
				break
			}
		}
	}
	lh.loadRequestQueue = queue
	return readyToLoad
}

// The mininum duration to sleep to load the next item in queue
// or Infinity if nothing in queue
func (lh *LoadHandler) findMinimumNextLoadDuration() time.Duration {

	lh.loadLock.Lock()
	defer lh.loadLock.Unlock()

	queue := lh.loadRequestQueue
	minNextLoadTime := lh.findMinimumNextLoadTime(queue)
	if minNextLoadTime == nil {
		return Infinity
	}
	duration := minNextLoadTime.Sub(time.Now())
	return duration
}

// The mininum next load time of any item.  Used to determine sleep time
// Can return nil if queue is empty
func (lh *LoadHandler) findMinimumNextLoadTime(queue []*LoadRequest) *time.Time {
	var minNextLoadTime *time.Time
	for _, loadRequest := range queue {
		if minNextLoadTime == nil || minNextLoadTime.After(*loadRequest.loadAfter) {
			minNextLoadTime = loadRequest.loadAfter
		}
	}
	return minNextLoadTime
}

func (lh *LoadHandler) load(loadRequest *LoadRequest) {

	if loadRequest == nil {
		return
	}

	toplog.Info("Metadata loader - Load data NOW. loadRequest: %+v", loadRequest)

	metadataHandler := metadataHandlerMap[loadRequest.dataType]
	if metadataHandler == nil {
		toplog.Error("Metadata loader - MetadataLoadHandler not found: %v", loadRequest)
		return
	}
	err := metadataHandler.MetadataLoadMethod(loadRequest.dataKey)
	if err != nil {
		// re-queue to load later if not at max load attempts
		loadRequest.loadAttempts = loadRequest.loadAttempts + 1
		if loadRequest.loadAttempts < MaxLoadAttempts {
			loadAfter := time.Now().Add(WaitToReloadOnErrorDuration)
			loadRequest.loadAfter = &loadAfter
			lh.RequestLoad(loadRequest)
		}
	}

}
