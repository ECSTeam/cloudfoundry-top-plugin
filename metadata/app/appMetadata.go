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
	//*common.BaseMetadataItem
	*common.Metadata
	*App
}

func NewAppMetadata(app App) *AppMetadata {
	appMetadata := &AppMetadata{}
	appMetadata.Metadata = common.NewMetadata()
	appMetadata.App = &app
	return appMetadata
}

func NewAppMetadataById(appId string) *AppMetadata {
	return NewAppMetadata(App{EntityCommon: common.EntityCommon{Guid: appId}, Name: appId})
}

func (metadataItem *AppMetadata) GetName() string {
	return metadataItem.Name
}
