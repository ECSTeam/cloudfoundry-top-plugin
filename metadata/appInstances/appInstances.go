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
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
)

// ****************************************************************
// The following are used used calling API: /v2/apps/APP_GUID/instances
// ****************************************************************
type AppInstances struct {
	//*common.BaseMetadataItem
	*common.Metadata
	*common.EntityCommon
	Data map[string]*AppInstance
	Guid string
	Name string
}

func NewAppInstances(appId string) *AppInstances {
	data := make(map[string]*AppInstance)
	appInst := NewAppInstancesWithData(appId, data)
	return appInst
}

func NewAppInstancesWithData(appId string, data map[string]*AppInstance) *AppInstances {
	appInst := &AppInstances{Guid: appId, Name: appId, Data: data}
	appInst.Metadata = common.NewMetadata()
	return appInst
}

func (metadataItem *AppInstances) GetGuid() string {
	return metadataItem.Guid
}

func (metadataItem *AppInstances) GetName() string {
	return metadataItem.Name
}

type AppInstance struct {
	Details   string     `json:"details"`
	Since     float64    `json:"since"`
	State     string     `json:"state"`
	Uptime    int64      `json:"uptime"`
	StartTime *time.Time // This will be populated on post-processing of response

}
