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

package space

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const UnknownName = "unknown"

type SpaceResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
	Resources []SpaceResource `json:"resources"`
}

type SpaceResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity Space       `json:"entity"`
}

type Space struct {
	common.EntityCommon
	//Guid                 string `json:"guid"`
	Name                 string `json:"name"`
	OrgGuid              string `json:"organization_guid"`
	QuotaGuid            string `json:"space_quota_definition_guid"`
	IsolationSegmentGuid string `json:"isolation_segment_guid"`
	Managers_url         string `json:"managers_url"`
	Auditors_url         string `json:"auditors_url"`
	Developers_url       string `json:"developers_url"`
}

/*
var (
	spacesMetadataCache []Space
)

func All() []Space {
	return spacesMetadataCache
}

func FindSpaceMetadata(spaceGuid string) Space {
	for _, space := range spacesMetadataCache {
		if space.Guid == spaceGuid {
			return space
		}
	}
	return Space{}
}

func FindSpaceName(spaceGuid string) string {
	spaceMetadata := FindSpaceMetadata(spaceGuid)
	spaceName := spaceMetadata.Name
	if spaceName == "" {
		spaceName = UnknownName
	}
	return spaceName
}

func LoadSpaceCache(cliConnection plugin.CliConnection) {
	data, err := getSpaceMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** space metadata error: %v", err.Error())
		return
	}
	spacesMetadataCache = data
}

func getSpaceMetadata(cliConnection plugin.CliConnection) ([]Space, error) {

	url := "/v2/spaces"
	metadata := []Space{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response SpaceResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			if item.Entity.IsolationSegmentGuid == "" {
				item.Entity.IsolationSegmentGuid = isolationSegment.DefaultIsolationSegmentGuid
			}
			metadata = append(metadata, item.Entity)
		}
		return response, response.NextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}
*/
