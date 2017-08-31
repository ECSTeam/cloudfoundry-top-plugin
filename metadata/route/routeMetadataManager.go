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

package route

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type RouteMetadataManager struct {
	*common.CommonV2ResponseManager
	internalRoutesMetadataCache []*RouteMetadata
	// Key: routeId, value: list of AppId
	appsForRouteCache map[string][]string
}

func NewRouteMetadataManager(mdGlobalManager common.MdGlobalManagerInterface) *RouteMetadataManager {
	url := "/v2/routes"
	mdMgr := &RouteMetadataManager{}
	mdMgr.CommonV2ResponseManager = common.NewCommonV2ResponseManager(mdGlobalManager, url, mdMgr, true)
	mdMgr.appsForRouteCache = make(map[string][]string)
	return mdMgr
}

func (mdMgr *RouteMetadataManager) CreateInternalGeneratedRoute(hostName string, pathName string, domainGuid string, port int) *RouteMetadata {
	route := &Route{
		EntityCommon:      common.EntityCommon{Guid: util.Pseudo_uuid()},
		Host:              hostName,
		Path:              pathName,
		DomainGuid:        domainGuid,
		Port:              port,
		InternalGenerated: true,
	}
	routeMd := NewRouteMetadata(*route)
	mdMgr.internalRoutesMetadataCache = append(mdMgr.internalRoutesMetadataCache, routeMd)
	return routeMd
}

func (mdMgr *RouteMetadataManager) FindItem(guid string) *RouteMetadata {
	foundRouteMd := mdMgr.FindItemInternal(guid, false, false)
	if foundRouteMd == nil {
		for _, route := range mdMgr.internalRoutesMetadataCache {
			if route.Guid == guid {
				return route
			}
		}
		return mdMgr.NewItemById(guid).(*RouteMetadata)
	}
	return foundRouteMd.(*RouteMetadata)
}

func (mdMgr *RouteMetadataManager) GetAll() []*RouteMetadata {
	// TODO: Need to use parent lock
	//mdMgr.mu.Lock()
	//defer mdMgr.mu.Unlock()
	appsMetadataArray := []*RouteMetadata{}
	for _, appMetadata := range mdMgr.MetadataMap {
		appsMetadataArray = append(appsMetadataArray, appMetadata.(*RouteMetadata))
	}
	return appsMetadataArray
}

func (mdMgr *RouteMetadataManager) NewItemById(guid string) common.IMetadata {
	return NewRouteMetadataById(guid)
}

func (mdMgr *RouteMetadataManager) CreateResponseObject() common.IResponse {
	return &RouteResponse{}
}

func (mdMgr *RouteMetadataManager) CreateResourceObject() common.IResource {
	return &RouteResource{}
}

func (mdMgr *RouteMetadataManager) CreateMetadataEntityObject(guid string) common.IMetadata {
	return NewRouteMetadataById(guid)
}

func (mdMgr *RouteMetadataManager) ProcessResponse(response common.IResponse, metadataArray []common.IMetadata) []common.IMetadata {
	resp := response.(*RouteResponse)
	for _, item := range resp.Resources {
		itemMd := mdMgr.ProcessResource(&item)
		metadataArray = append(metadataArray, itemMd)
	}
	return metadataArray
}

func (mdMgr *RouteMetadataManager) ProcessResource(resource common.IResource) common.IMetadata {
	resourceType := resource.(*RouteResource)
	resourceType.Entity.Guid = resourceType.Meta.Guid
	metadata := NewRouteMetadata(resourceType.Entity)
	return metadata
}

func (mdMgr *RouteMetadataManager) FindAppIdsForRouteMetadata(routeGuid string) []string {
	appIds := mdMgr.appsForRouteCache[routeGuid]
	if appIds == nil {
		// We stick an empty array in to prevent triggering go routine multiple times
		// TODO: Find a better way to do this
		// TODO: Need a way to do a callback / tickle when metadata is loaded so
		// caller who wanted the data can refresh screen (if still relevant)
		mdMgr.appsForRouteCache[routeGuid] = make([]string, 0)
		go mdMgr.LoadAppsForRouteCache(routeGuid)
	}
	return appIds
}

func (mdMgr *RouteMetadataManager) LoadAppsForRouteCache(routeId string) {
	appIds := mdMgr.getAppIdsForRoute(routeId)
	if appIds != nil {
		mdMgr.appsForRouteCache[routeId] = appIds
	}
}

func (mdMgr *RouteMetadataManager) getAppIdsForRoute(routeId string) []string {
	appList, err := mdMgr.getAppsForRoute(routeId)
	if err != nil {
		toplog.Warn("*** getAppsForRoute metadata error: %v", err.Error())
		return nil
	}
	appIdList := make([]string, len(appList))
	for i, app := range appList {
		appIdList[i] = app.GetGuid()
	}
	return appIdList
}

func (mdMgr *RouteMetadataManager) getAppsForRoute(routeId string) ([]common.IMetadata, error) {
	url := fmt.Sprintf("/v2/routes/%v/apps", routeId)
	toplog.Debug("getAppsForRoute url: %v", url)
	return mdMgr.GetMdGlobalManager().GetAppMetadataFromUrl(url)
}
