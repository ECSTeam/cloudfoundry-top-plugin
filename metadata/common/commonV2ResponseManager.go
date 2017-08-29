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

package common

import (
	"encoding/json"
	"time"

	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type V2MetadataManager interface {
	MetadataManager
	LoadAllItemsInternal() ([]IMetadata, error)

	CreateResponseObject() IResponse
	CreateResourceObject() IResource
	CreateMetadataEntityObject(guid string) IMetadata
	ProcessResponse(IResponse, []IMetadata) []IMetadata
	ProcessResource(resource IResource) IMetadata
	GetNextUrl(response IResponse) string
}

type CommonV2ResponseManager struct {
	*CommonMetadataManager
	mm V2MetadataManager

	autoFullLoadIfNotFound bool
	fullLoadCacheTime      time.Time
	loadInProgress         bool
}

func NewCommonV2ResponseManager(mdGlobalManager MdGlobalManagerInterface,
	url string,
	mm V2MetadataManager,
	autoFullLoadIfNotFound bool) *CommonV2ResponseManager {

	commonV2ResponseMgr := &CommonV2ResponseManager{mm: mm, autoFullLoadIfNotFound: autoFullLoadIfNotFound}
	commonV2ResponseMgr.CommonMetadataManager = NewCommonMetadataManager(mdGlobalManager, url, mm)
	return commonV2ResponseMgr
}

func (commonV2ResponseMgr *CommonV2ResponseManager) FindItemInternal(guid string, requestLoadIfNotFound bool, createEmptyObjectIfNotFound bool) IMetadata {

	requestLoadIfNotFound = requestLoadIfNotFound && !commonV2ResponseMgr.autoFullLoadIfNotFound
	item := commonV2ResponseMgr.CommonMetadataManager.FindItemInternal(guid, requestLoadIfNotFound, createEmptyObjectIfNotFound)

	if commonV2ResponseMgr.autoFullLoadIfNotFound && (item == nil || item.GetCacheTime() == nil) {
		commonV2ResponseMgr.LoadAllItemsAysnc()
	}
	return item
}

func (commonV2ResponseMgr *CommonV2ResponseManager) LoadAllItemsInternal() ([]IMetadata, error) {
	return commonV2ResponseMgr.GetMetadata()
}

func (commonV2ResponseMgr *CommonV2ResponseManager) LoadItemInternal(guid string) (IMetadata, error) {
	url := commonV2ResponseMgr.url + "/" + guid
	now := time.Now()

	outputStr, err := CallAPI(commonV2ResponseMgr.mdGlobalManager.GetCliConnection(), url)
	if err != nil {
		emptyApp := commonV2ResponseMgr.mm.NewItemById(guid)
		return emptyApp, err
	}
	outputBytes := []byte(outputStr)
	resource := commonV2ResponseMgr.mm.CreateResourceObject()
	err = json.Unmarshal(outputBytes, resource)
	if err != nil {
		emptyApp := commonV2ResponseMgr.mm.NewItemById(guid)
		return emptyApp, err
	}

	itemMetadata := commonV2ResponseMgr.mm.ProcessResource(resource)
	itemMetadata.SetCacheTime(&now)
	return itemMetadata, nil
}

func (commonV2ResponseMgr *CommonV2ResponseManager) LoadAllItems() {
	now := time.Now()

	metadataItemArray, err := commonV2ResponseMgr.LoadAllItemsInternal()
	if err != nil {
		toplog.Warn("*** app metadata error: %v", err.Error())
		return
	}

	metadataMap := make(map[string]IMetadata)
	for _, metadataItem := range metadataItemArray {
		//toplog.Debug("From Map - app id: %v name:%v", appMetadata.Guid, appMetadata.Name)
		metadataItem.SetCacheTime(&now)
		metadataMap[metadataItem.GetGuid()] = metadataItem
	}
	commonV2ResponseMgr.MetadataMap = metadataMap
}

func (commonV2ResponseMgr *CommonV2ResponseManager) LoadAllItemsAysnc() {

	commonV2ResponseMgr.mu.Lock()
	defer commonV2ResponseMgr.mu.Unlock()

	if commonV2ResponseMgr.loadInProgress {
		toplog.Debug("CommonV2ResponseManager.LoadAllItemsAysnc %v loadInProgress", commonV2ResponseMgr.url)
		return
	}

	commonV2ResponseMgr.loadInProgress = true
	loadAsync := func() {
		toplog.Debug("CommonV2ResponseManager.LoadAllItemsAysnc %v loadAsync thread started", commonV2ResponseMgr.url)
		commonV2ResponseMgr.LoadAllItems()
		toplog.Debug("CommonV2ResponseManager.LoadAllItemsAysnc %v loadAsync thread complete", commonV2ResponseMgr.url)
		commonV2ResponseMgr.loadInProgress = false
	}
	go loadAsync()
}

func (commonMgr *CommonMetadataManager) GetNextUrl(response IResponse) string {
	nextUrl, _ := GetStringValueByFieldName(response, "NextUrl")
	return nextUrl
}

func (commonV2ResponseMgr *CommonV2ResponseManager) GetMetadata() ([]IMetadata, error) {
	return commonV2ResponseMgr.GetMetadataFromUrl(commonV2ResponseMgr.GetUrl())
}

func (commonV2ResponseMgr *CommonV2ResponseManager) GetMetadataFromUrl(url string) ([]IMetadata, error) {
	metadataArray := []IMetadata{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		resp := commonV2ResponseMgr.mm.CreateResponseObject()
		err = json.Unmarshal(outputBytes, &resp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadataArray, "", err
		}
		metadataArray = commonV2ResponseMgr.mm.ProcessResponse(resp, metadataArray)

		nextUrl = commonV2ResponseMgr.mm.GetNextUrl(resp)
		//nextUrl, _ = GetStringValueByFieldName(resp, "NextUrl")
		return resp, nextUrl, nil
	}

	err := CallPagableAPI(commonV2ResponseMgr.mdGlobalManager.GetCliConnection(), url, handleRequest)

	return metadataArray, err

}

/*
func (commonV2ResponseMgr *CommonV2ResponseManager) GetMetadataV3FromUrl() ([]IMetadata, error) {

	url := commonV2ResponseMgr.GetUrl()
	metadataArray := []IMetadata{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		resp := commonV2ResponseMgr.mm.CreateResponseObject().(IResponseV3)
		err = json.Unmarshal(outputBytes, &resp)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadataArray, "", err
		}
		metadataArray = commonV2ResponseMgr.mm.ProcessResponse(resp, metadataArray)
		nextUrl = resp.GetPagination().Next.Href
		return resp, nextUrl, nil
	}

	err := CallPagableAPI(commonV2ResponseMgr.mdGlobalManager.GetCliConnection(), url, handleRequest)

	return metadataArray, err

}
*/
