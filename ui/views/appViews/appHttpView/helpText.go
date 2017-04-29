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

package appHttpView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText +
	helpView.HelpHeaderText +
	HelpColumnsText +
	HelpLocalViewKeybindings +
	helpView.HelpChildLevelDataViewKeybindings +
	helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**App HTTP(S) Response Info**

App HTTP(S) response info view shows all the HTTP and HTTPS responses
that have occured from the selected application. 
`

const HelpColumnsText = `
**HTTP(S) Response Columns:**

  METHOD - The HTTP method used for the request
  CODE - The HTTP response code.
  LAST_RESPONSE - Last time this METHOD+CODE combination occured. 
  COUNT - Number of responses that have occured for this METHOD+CODE
          combination.
`

const HelpLocalViewKeybindings = `
`
