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

type StackResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
	Resources []StackResource `json:"resources"`
}

type StackResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Stack `json:"entity"`
}

type Stack struct {
	Guid        string `json:"guid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	stacksMetadataCache []Stack
)

func AllStacks() []Stack {
	return stacksMetadataCache
}

func FindStackMetadata(stackGuid string) Stack {
	for _, stack := range stacksMetadataCache {
		if stack.Guid == stackGuid {
			return stack
		}
	}
	return Stack{Guid: stackGuid, Name: stackGuid}
}

func LoadStackCache(cliConnection plugin.CliConnection) {
	data, err := getStackMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** stack metadata error: %v", err.Error()))
		return
	}
	stacksMetadataCache = data
}

func getStackMetadata(cliConnection plugin.CliConnection) ([]Stack, error) {

	url := "/v2/stacks"
	metadata := []Stack{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response StackResponse
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

	err := callPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}
