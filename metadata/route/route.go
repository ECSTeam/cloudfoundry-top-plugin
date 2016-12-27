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

package route

import (
	"encoding/json"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type RouteResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
	Resources []RouteResource `json:"resources"`
}

type RouteResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity Route       `json:"entity"`
}

type Route struct {
	Guid                string `json:"guid"`
	Host                string `json:"host"`
	Path                string `json:"path"`
	DomainGuid          string `json:"domain_guid"`
	SpaceGuid           string `json:"space_guid"`
	ServiceInstanceGuid string `json:"service_instance_guid"`
	Port                int    `json:"port"`
	InternalGenerated   bool
}

func CreateInternalGeneratedRoute(hostName string, pathName string, domainGuid string, port int) *Route {
	r := &Route{
		Guid:              util.Pseudo_uuid(),
		Host:              hostName,
		Path:              pathName,
		DomainGuid:        domainGuid,
		Port:              port,
		InternalGenerated: true,
	}
	internalRoutesMetadataCache = append(internalRoutesMetadataCache, r)
	return r
}

var (
	routesMetadataCache         []*Route
	internalRoutesMetadataCache []*Route
)

func AllRoutes() []*Route {
	return routesMetadataCache
}

func FindRouteMetadata(routeGuid string) *Route {
	for _, route := range routesMetadataCache {
		if route.Guid == routeGuid {
			return route
		}
	}
	for _, route := range internalRoutesMetadataCache {
		if route.Guid == routeGuid {
			return route
		}
	}
	return &Route{Guid: routeGuid}
}

func LoadRouteCache(cliConnection plugin.CliConnection) {
	data, err := getRouteMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** route metadata error: %v", err.Error())
		return
	}
	routesMetadataCache = data
}

func getRouteMetadata(cliConnection plugin.CliConnection) ([]*Route, error) {

	url := "/v2/routes"
	metadata := []*Route{}

	toplog.Debug("Route>>getRouteMetadata start")

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response RouteResponse
		err := json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			//itemMetadata := NewRouteMetadata(item.Entity)
			entity := item.Entity
			metadata = append(metadata, &entity)
		}
		return response, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	toplog.Debug("Route>>getRouteMetadata complete - loaded: %v items", len(metadata))

	return metadata, err

}
