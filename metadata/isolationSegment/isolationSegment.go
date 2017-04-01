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

package isolationSegment

import (
	"encoding/json"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

const UnknownName = "unknown"
const SharedIsolationSegmentName = "shared"

type Link struct {
	Href string `json:"href"`
}

type Pagination struct {
	Count int  `json:"total_results"`
	Pages int  `json:"total_pages"`
	Next  Link `json:"next"`
}

type IsolationSegmentResponse struct {
	Pagination Pagination         `json:"pagination"`
	Resources  []IsolationSegment `json:"resources"`
}

type IsolationSegment struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

var (
	sharedIsolationSegment        = IsolationSegment{Name: SharedIsolationSegmentName}
	isolationSegmentMetadataCache []IsolationSegment
)

func All() []IsolationSegment {
	return isolationSegmentMetadataCache
}

func GetDefault() IsolationSegment {
	return sharedIsolationSegment
}

func FindMetadata(guid string) IsolationSegment {
	if guid == "" {
		return sharedIsolationSegment
	}
	for _, isoSeg := range isolationSegmentMetadataCache {
		if isoSeg.Guid == guid {
			return isoSeg
		}
	}
	return IsolationSegment{Guid: guid, Name: guid}
}

func FindName(guid string) string {
	metadata := FindMetadata(guid)
	name := metadata.Name
	if name == "" {
		name = UnknownName
	}
	return name
}

func FindMetadataByName(name string) IsolationSegment {
	if name == "" {
		return sharedIsolationSegment
	}
	for _, isoSeg := range isolationSegmentMetadataCache {
		if isoSeg.Name == name {
			return isoSeg
		}
	}
	return IsolationSegment{Name: name}
}

func LoadCache(cliConnection plugin.CliConnection) {
	data, err := getMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** isolationSegment metadata error: %v", err.Error())
		return
	}

	//toplog.Info("isolation segments: %+v", data)
	isolationSegmentMetadataCache = data
	sharedIsolationSegment = FindMetadataByName(SharedIsolationSegmentName)

}

func getMetadata(cliConnection plugin.CliConnection) ([]IsolationSegment, error) {

	url := "/v3/isolation_segments"
	metadata := []IsolationSegment{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response IsolationSegmentResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			metadata = append(metadata, item)
		}
		return response, response.Pagination.Next.Href, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}
