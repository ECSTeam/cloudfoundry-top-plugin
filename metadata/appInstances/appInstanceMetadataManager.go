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

package appInstances

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type AppInstanceMetadataManager struct {
	*common.CommonMetadataManager
	mu sync.Mutex
}

func NewAppInstanceMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *AppInstanceMetadataManager {
	url := "/v2/apps"
	mdMgr := &AppInstanceMetadataManager{}
	mdMgr.CommonMetadataManager = common.NewCommonMetadataManager(mdGlobalManager, common.APP_INST, url, mdMgr)
	return mdMgr
}

func (mdMgr *AppInstanceMetadataManager) MinimumReloadDuration() time.Duration {
	return time.Millisecond * 1000
}

func (mdMgr *AppInstanceMetadataManager) FindItem(appId string) *AppInstances {
	item := mdMgr.FindItemInternal(appId, false, false)
	if item != nil {
		return item.(*AppInstances)
	}
	return nil
}

func (mdMgr *AppInstanceMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewAppInstances(guid)
}

func (mdMgr *AppInstanceMetadataManager) LoadItemInternal(guid string) (common.IMetadata, error) {

	url := mdMgr.GetUrl() + "/" + guid + "/instances"

	output, err := common.CallAPI(mdMgr.GetMdGlobalManager().GetCliConnection(), url)
	if err != nil {
		return nil, err
	}

	if strings.Contains(output, "error_code") {
		// "Instances error: Request failed for app: cf-nodejs as the app is in stopped state."
		if strings.Contains(output, "220001") {
			// This error is OK
			return NewAppInstances(guid), nil
		} else {
			errMsg := fmt.Sprintf("Error from API call: %v", output)
			return nil, errors.New(errMsg)
		}
	}

	response := make(map[string]*AppInstance)
	outputBytes := []byte(output)
	err = json.Unmarshal(outputBytes, &response)
	if err != nil {
		toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
		return NewAppInstances(guid), err
	}

	// Set the startTime relative to now and uptime of the container
	for _, stat := range response {
		// Ignore "uptime" field if container is in state DOWN
		if stat.State == "DOWN" {
			stat.Uptime = 0
		} else {
			startTime := time.Unix(int64(stat.Since), 0)
			stat.StartTime = &startTime
		}
	}

	instances := NewAppInstancesWithData(guid, response)
	return instances, nil

}
