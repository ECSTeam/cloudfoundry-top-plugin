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

package spaceQuota

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const UnknownName = "unknown"

type SpaceQuotaResponse struct {
	Count     int                  `json:"total_results"`
	Pages     int                  `json:"total_pages"`
	Resources []SpaceQuotaResource `json:"resources"`
}

type SpaceQuotaResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity SpaceQuota  `json:"entity"`
}

type SpaceQuota struct {
	//Guid                    string `json:"guid"`
	common.EntityCommon
	Name                    string `json:"name"`
	OrganizationGuid        string `json:"organization_guid"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	MemoryLimit             int    `json:"memory_limit"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
}
