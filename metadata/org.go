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

package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type Org struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

var (
	orgsMetadataCache []Org
)

func FindOrgMetadata(orgGuid string) Org {
	for _, org := range orgsMetadataCache {
		if org.Guid == orgGuid {
			return org
		}
	}
	return Org{}
}

func LoadOrgCache(cliConnection plugin.CliConnection) {
	data, err := getOrgMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** org metadata error: %v", err.Error()))
		return
	}
	orgsMetadataCache = data
}

func getOrgMetadata(cliConnection plugin.CliConnection) ([]Org, error) {

	url := "/v2/organizations"
	metadata := []Org{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response OrgResponse
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

	callAPI(cliConnection, url, handleRequest)

	return metadata, nil

}
