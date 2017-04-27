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

package appHttpView

import (
	"fmt"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata/eventApp"
)

type DisplayHttpInfo struct {
	*eventApp.HttpInfo
	LastAcivityFormatted string
	key                  string
}

func NewDisplayHttpInfo(info *eventApp.HttpInfo) *DisplayHttpInfo {
	displayInfo := &DisplayHttpInfo{}
	displayInfo.HttpInfo = info
	return displayInfo
}

func (dhi *DisplayHttpInfo) Id() string {
	if dhi.key == "" {
		dhi.key = fmt.Sprintf("%v-%v", dhi.HttpMethod, dhi.HttpStatusCode)
	}
	return dhi.key
}
