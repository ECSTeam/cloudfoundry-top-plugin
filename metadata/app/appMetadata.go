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

package app

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const MEGABYTE = (1024 * 1024)

type AppMetadata struct {
	*common.BaseMetadataItem
	*App
}

func NewAppMetadata(app App) *AppMetadata {
	appMetadata := &AppMetadata{}
	appMetadata.BaseMetadataItem = common.NewBaseMetadataItem()
	appMetadata.App = &app
	return appMetadata
}

func NewAppMetadataById(appId string) *AppMetadata {
	return NewAppMetadata(App{Guid: appId, Name: appId})
}

func (metadataItem *AppMetadata) GetGuid() string {
	return metadataItem.Guid
}

func (metadataItem *AppMetadata) GetName() string {
	return metadataItem.Name
}
