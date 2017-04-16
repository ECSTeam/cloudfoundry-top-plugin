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

package crashData

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const UnknownName = "unknown"

type EventDataResponse struct {
	Count     int                 `json:"total_results"`
	Pages     int                 `json:"total_pages"`
	NextUrl   string              `json:"next_url"`
	Resources []EventDataResource `json:"resources"`
}

type EventDataResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity EventData   `json:"entity"`
}

type EventData struct {
	common.EntityCommon
	Type              string                 `json:"type"`
	Actor             string                 `json:"actor"`
	Actor_type        string                 `json:"actor_type"`
	Actor_name        string                 `json:"actor_name"`
	Actor_username    string                 `json:"actor_username"`
	Actee             string                 `json:"actee"`
	Actee_type        string                 `json:"actee_type"`
	Actee_name        string                 `json:"actee_name"`
	Timestamp         string                 `json:"timestamp"`
	Space_guid        string                 `json:"space_guid"`
	Organization_guid string                 `json:"organization_guid"`
	Metadata          EventDataMetadataField `json:"metadata"`
}

type EventDataMetadataField struct {
	Instance         string `json:"instance"`
	Index            int    `json:"index"`
	Exit_description string `json:"exit_description"`
	Reason           string `json:"reason"`
}
