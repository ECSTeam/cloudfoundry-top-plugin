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

type AppResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []AppResource `json:"resources"`
}

type AppResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity App         `json:"entity"`
}

type App struct {
	common.EntityCommon

	//Guid      string `json:"guid"`
	Name      string `json:"name,omitempty"`
	SpaceGuid string `json:"space_guid,omitempty"`

	StackGuid   string  `json:"stack_guid,omitempty"`
	MemoryMB    float64 `json:"memory,omitempty"`
	DiskQuotaMB float64 `json:"disk_quota,omitempty"`

	Environment map[string]interface{} `json:"environment_json,omitempty"`
	Instances   float64                `json:"instances,omitempty"`
	State       string                 `json:"state,omitempty"`
	EnableSsh   bool                   `json:"enable_ssh,omitempty"`

	PackageState        string `json:"package_state,omitempty"`
	StagingFailedReason string `json:"staging_failed_reason,omitempty"`
	StagingFailedDesc   string `json:"staging_failed_description,omitempty"`
	DetectedStartCmd    string `json:"detected_start_command,omitempty"`
	DockerImage         string `json:"docker_image,omitempty"`
	//DockerCredentials string  `json:"docker_credentials_json,omitempty"`
	//audit.app.create event fields
	Console           bool   `json:"console,omitempty"`
	Buildpack         string `json:"buildpack,omitempty"`
	DetectedBuildpack string `json:"detected_buildpack,omitempty"`

	HealthcheckType    string  `json:"health_check_type,omitempty"`
	HealthcheckTimeout float64 `json:"health_check_timeout,omitempty"`
	Production         bool    `json:"production,omitempty"`
	//app.crash event fields
	//Index           float64 `json:"index,omitempty"`
	//ExitStatus      string  `json:"exit_status,omitempty"`
	//ExitDescription string  `json:"exit_description,omitempty"`
	//ExitReason      string  `json:"reason,omitempty"`
	// "package_updated_at": "2016-11-15T19:56:52Z",
	PackageUpdatedAt string `json:"package_updated_at,omitempty"`
}
