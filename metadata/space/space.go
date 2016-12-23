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

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

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
	Guid    string `json:"guid"`
	Name    string `json:"name"`
	OrgGuid string `json:"organization_guid"`
	OrgName string
}

var (
	spacesMetadataCache []Space
)

func FindSpaceMetadata(spaceGuid string) Space {
	for _, space := range spacesMetadataCache {
		if space.Guid == spaceGuid {
			return space
		}
	}
	return Space{}
}

func LoadSpaceCache(cliConnection plugin.CliConnection) {
	data, err := getSpaceMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** space metadata error: %v", err.Error()))
		return
	}
	spacesMetadataCache = data
}

func getSpaceMetadata(cliConnection plugin.CliConnection) ([]Space, error) {

	url := "/v2/spaces"
	metadata := []Space{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response SpaceResponse
		err := json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn(fmt.Sprintf("*** %v unmarshal parsing output: %v", url, string(outputBytes[:])))
			return metadata, err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			metadata = append(metadata, item.Entity)
		}
		return response, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}
