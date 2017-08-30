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

package domain

import "github.com/ecsteam/cloudfoundry-top-plugin/metadata/common"

type DomainResponse struct {
	Count     int              `json:"total_results"`
	Pages     int              `json:"total_pages"`
	NextUrl   string           `json:"next_url"`
	Resources []DomainResource `json:"resources"`
}

type DomainResource struct {
	Meta   common.Meta `json:"metadata"`
	Entity Domain      `json:"entity"`
}

type Domain struct {
	common.EntityCommon
	//Guid                   string `json:"guid"`
	Name                   string `json:"name"`
	RouterGroupGuid        string `json:"router_group_guid"`
	RouterGroupType        string `json:"router_group_type"`
	OwningOrganizationGuid string `json:"owning_organization_guid"`

	SharedDomain bool
}

/*
var (
	domainsMetadataCache []*Domain
)


func AllDomains() []*Domain {
	return domainsMetadataCache
}

func FindDomainMetadata(domainGuid string) *Domain {
	for _, domain := range domainsMetadataCache {
		if domain.Guid == domainGuid {
			return domain
		}
	}
	return &Domain{Guid: domainGuid}
}

func FindDomainMetadataByName(domainName string) *Domain {
	for _, domain := range domainsMetadataCache {
		if domain.Name == domainName {
			return domain
		}
	}
	return nil
}

func AddDomainMetadata(domainName string) *Domain {
	domain := &Domain{Guid: util.Pseudo_uuid(), Name: domainName}
	domainsMetadataCache = append(domainsMetadataCache, domain)
	return domain
}

func LoadDomainCache(cliConnection plugin.CliConnection) {
	sharedDomains, err := getDomainMetadata(cliConnection, "/v2/shared_domains")
	if err != nil {
		toplog.Warn("*** shared_domains metadata error: %v", err.Error())
		return
	}
	for _, domain := range sharedDomains {
		domain.SharedDomain = true
	}

	privateDomains, err := getDomainMetadata(cliConnection, "/v2/private_domains")
	if err != nil {
		toplog.Warn("*** private_domains metadata error: %v", err.Error())
		return
	}

	data := append(sharedDomains, privateDomains...)
	toplog.Debug("Domain>>LoadDomainCache total items loaded: %v", len(data))
	domainsMetadataCache = data
}

func getDomainMetadata(cliConnection plugin.CliConnection, url string) ([]*Domain, error) {

	metadata := []*Domain{}

	toplog.Debug("Domain>>getDomainMetadata %v start", url)

	handleRequest := func(outputBytes []byte) (data interface{}, nextUrl string, err error) {
		var response DomainResponse
		err = json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn("*** %v unmarshal parsing output: %v", url, string(outputBytes[:]))
			return metadata, "", err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			//itemMetadata := NewDomainMetadata(item.Entity)
			entity := item.Entity
			metadata = append(metadata, &entity)
		}
		return response, response.NextUrl, nil
	}

	err := common.CallPagableAPI(cliConnection, url, handleRequest)

	toplog.Debug("Domain>>getDomainMetadata %v complete - loaded: %v items", url, len(metadata))

	return metadata, err

}
*/
