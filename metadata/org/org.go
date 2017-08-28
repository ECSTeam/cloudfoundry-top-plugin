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

package org

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

const UnknownName = "unknown"

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity Org         `json:"entity"`
}

type Org struct {
	common.EntityCommon
	//Guid                 string `json:"guid"`
	Name                 string `json:"name"`
	QuotaGuid            string `json:"quota_definition_guid"`
	Status               string `json:"status"`
	Domains_url          string `json:"domains_url"`
	Private_domains_url  string `json:"private_domains_url"`
	Users_url            string `json:"users_url"`
	Managers_url         string `json:"managers_url"`
	Auditors_url         string `json:"auditors_url"`
	Billing_managers_url string `json:"billing_managers_url"`
}

/*
var (
	orgsMetadataCache []Org
)

func All() []Org {
	return orgsMetadataCache
}

func FindOrgMetadata(orgGuid string) Org {
	for _, org := range orgsMetadataCache {
		if org.Guid == orgGuid {
			return org
		}
	}
	return Org{}
}

func FindOrgNameBySpaceGuid(spaceGuid string) string {
	_, orgName := FindBySpaceGuid(spaceGuid)
	return orgName
}

func FindBySpaceGuid(spaceGuid string) (orgId string, orgName string) {
	spaceMetadata := space.FindSpaceMetadata(spaceGuid)
	orgId = spaceMetadata.OrgGuid
	orgMetadata := FindOrgMetadata(orgId)
	orgName = orgMetadata.Name
	//toplog.Info("Lookup name for org via space guid: %v found name:[%v]", spaceGuid, orgName)
	if orgName == "" {
		orgName = UnknownName
	}
	return orgId, orgName
}

func LoadOrgCache(cliConnection plugin.CliConnection) {
	data, err := getOrgMetadata(cliConnection)
	if err != nil {
		toplog.Warn("*** org metadata error: %v", err.Error())
		return
	}
	orgsMetadataCache = data
}

func getOrgMetadata(cliConnection plugin.CliConnection) ([]Org, error) {

	url := "/v2/organizations"
	metadata := []Org{}

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response OrgResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			metadata = append(metadata, item.Entity)
		}
		return response, response.NextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	return metadata, err

}

*/
