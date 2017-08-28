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

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const SharedIsolationSegmentName = "shared"
const DefaultIsolationSegmentGuid = "-1"
const UnknownIsolationSegmentGuid = ""
const UnknownIsolationSegmentName = "unknown"

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

func (resp *IsolationSegmentResponse) GetPagination() Pagination {
	return resp.Pagination
}

type IsolationSegment struct {
	common.EntityCommon
	//Guid string `json:"guid"`
	Name string `json:"name"`
}

var (
	DefaultIsolationSegment = NewIsolationSegmentMetadata(IsolationSegment{EntityCommon: common.EntityCommon{Guid: DefaultIsolationSegmentGuid}, Name: "default"})
	UnknownIsolationSegment = NewIsolationSegmentMetadata(IsolationSegment{EntityCommon: common.EntityCommon{Guid: UnknownIsolationSegmentGuid}, Name: UnknownIsolationSegmentName})
	SharedIsolationSegment  *IsolationSegmentMetadata //= NewIsolationSegmentMetadata(IsolationSegment{})
)

func GetDefault() *IsolationSegmentMetadata {
	return SharedIsolationSegment
}

/*
func All() []*IsolationSegment {
	return isolationSegmentMetadataCache
}

func GetDefault() *IsolationSegment {
	return SharedIsolationSegment
}

func FindMetadata(guid string) *IsolationSegment {
	if guid == "" {
		return &IsolationSegment{Name: "unknown"}
	}
	if guid == DefaultIsolationSegmentGuid {
		return SharedIsolationSegment
	}
	for _, isoSeg := range isolationSegmentMetadataCache {
		if isoSeg.Guid == guid {
			return isoSeg
		}
	}
	return &IsolationSegment{Guid: guid, Name: guid}
}

func FindName(guid string) string {
	metadata := FindMetadata(guid)
	name := metadata.Name
	if name == "" {
		name = UnknownIsolationSegmentName
	}
	return name
}

func FindMetadataByName(name string) *IsolationSegment {
	if name == "" {
		return SharedIsolationSegment
	}
	for _, isoSeg := range isolationSegmentMetadataCache {
		if isoSeg.Name == name {
			return isoSeg
		}
	}
	return &IsolationSegment{Name: name}
}

func LoadCache(cliConnection plugin.CliConnection) {
	data, err := getMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** isolationSegment metadata error: %v", err.Error())
		return
	}

	//toplog.Info("isolation segments: %+v", data)
	isolationSegmentMetadataCache = data
	SharedIsolationSegment = FindMetadataByName(SharedIsolationSegmentName)
}

func getMetadata(cliConnection plugin.CliConnection) ([]*IsolationSegment, error) {

	url := "/v3/isolation_segments"
	metadata := []*IsolationSegment{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response IsolationSegmentResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			item := &IsolationSegment{Guid: item.Guid, Name: item.Name}
			metadata = append(metadata, item)
		}
		return response, response.Pagination.Next.Href, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}
*/
