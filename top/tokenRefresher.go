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

package top

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

type TokenRefresher struct {
	cliConnection    plugin.CliConnection
	nozzleInstanceId int
}

func NewTokenRefresher(cliConnection plugin.CliConnection, nozzleInstanceId int) *TokenRefresher {
	return &TokenRefresher{
		cliConnection:    cliConnection,
		nozzleInstanceId: nozzleInstanceId,
	}
}

func (tr *TokenRefresher) RefreshAuthToken() (string, error) {
	toplog.Info("Nozzle #%v - RefreshAuthToken called", tr.nozzleInstanceId)
	token, err := tr.cliConnection.AccessToken()
	if err != nil {
		toplog.Error("Nozzle #%v - RefreshAuthToken failed: %v", tr.nozzleInstanceId, err)
		return "", err
	}
	toplog.Info("Nozzle #%v - RefreshAuthToken complete with new token: %v", tr.nozzleInstanceId, token)
	return token, nil
}
