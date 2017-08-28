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

import "time"

type Link struct {
	Href string `json:"href"`
}

type Pagination struct {
	Count int  `json:"total_results"`
	Pages int  `json:"total_pages"`
	Next  Link `json:"next"`
}

type IResponse interface {
	//Count() int
	//Pages() int
	//Resources() []IResource
}

type IResponseV3 interface {
	GetPagination() Pagination
}

type IResource interface {
	//Guid() string
	//Meta() IResourceMetadata
	//Entity() IEntity
}

type IResourceMetadata interface {
	//GetGuid() string
	//SetGuid(string)
	//SetCreatedAt(string)
	//SetUpdatedAt(string)
}

type IEntity interface {
	GetGuid() string
	SetGuid(string)
	GetName() string
}

type IMetadata interface {
	IEntity
	SetCacheTime(*time.Time)
	GetCacheTime() *time.Time
}

type Metadata struct {
	cacheTime *time.Time
}

func NewMetadata() *Metadata {
	item := &Metadata{}
	now := time.Now()
	item.cacheTime = &now
	return item
}

func (md *Metadata) SetCacheTime(time *time.Time) {
	md.cacheTime = time
}

func (md *Metadata) GetCacheTime() *time.Time {
	return md.cacheTime
}
