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

type DataType string

const (
	APP            DataType = "APP"
	APP_INST                = "APP_INST"
	APP_STATS               = "APP_STATS"
	SPACE                   = "SPACE"
	ORG                     = "ORG"
	DOMAIN_PRIVATE          = "DOMAIN_PRIVATE"
	DOMAIN_SHARED           = "DOMAIN_SHARED"
	ISO_SEG                 = "ISO_SEG"
	ORG_QUOTA               = "ORG_QUOTA"
	SPACE_QUOTA             = "SPACE_QUOTA"
	ROUTE                   = "ROUTE"
	STACK                   = "STACK"
	EVENTS_CRASH            = "EVENTS_CRASH"
)
