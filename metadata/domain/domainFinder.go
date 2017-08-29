// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.Domain/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

type DomainFinder struct {
	sharedMdMgr  *DomainSharedMetadataManager
	privateMdMgr *DomainPrivateMetadataManager
}

func NewDomainFinder(sharedMdMgr *DomainSharedMetadataManager, privateMdMgr *DomainPrivateMetadataManager) *DomainFinder {
	return &DomainFinder{sharedMdMgr: sharedMdMgr, privateMdMgr: privateMdMgr}
}

func (mdFinder *DomainFinder) FindDomainMetadata(guid string) *DomainMetadata {
	domainMd := mdFinder.findDomainSharedMetadataByGuid(guid)
	if domainMd == nil {
		domainMd = mdFinder.findDomainPrivateMetadataByGuid(guid)
	}
	return domainMd
}

func (mdFinder *DomainFinder) findDomainSharedMetadataByGuid(guid string) *DomainMetadata {
	return mdFinder.findDomainMetadataByGuidInArray(mdFinder.sharedMdMgr.GetAll(), guid)
}

func (mdFinder *DomainFinder) findDomainPrivateMetadataByGuid(guid string) *DomainMetadata {
	return mdFinder.findDomainMetadataByGuidInArray(mdFinder.privateMdMgr.GetAll(), guid)
}

func (mdFinder *DomainFinder) findDomainMetadataByGuidInArray(domains []*DomainMetadata, guid string) *DomainMetadata {
	for _, domain := range domains {
		if domain.Guid == guid {
			return domain
		}
	}
	return nil
}

func (mdFinder *DomainFinder) FindDomainMetadataByName(name string) *DomainMetadata {
	domainMd := mdFinder.findDomainSharedMetadataByName(name)
	if domainMd == nil {
		domainMd = mdFinder.findDomainPrivateMetadataByName(name)
	}
	return domainMd
}

func (mdFinder *DomainFinder) findDomainSharedMetadataByName(name string) *DomainMetadata {
	return mdFinder.findDomainMetadataByNameInArray(mdFinder.sharedMdMgr.GetAll(), name)
}

func (mdFinder *DomainFinder) findDomainPrivateMetadataByName(name string) *DomainMetadata {
	return mdFinder.findDomainMetadataByNameInArray(mdFinder.privateMdMgr.GetAll(), name)
}

func (mdFinder *DomainFinder) findDomainMetadataByNameInArray(domains []*DomainMetadata, name string) *DomainMetadata {
	for _, domain := range domains {
		if domain.Name == name {
			return domain
		}
	}
	return nil
}
