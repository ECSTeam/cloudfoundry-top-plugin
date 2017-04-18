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

package appCrashView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText +
	helpView.HelpHeaderText +
	HelpColumnsText +
	HelpLocalViewKeybindings +
	helpView.HelpChildLevelDataViewKeybindings +
	helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**App Container CRASH View**

App Container CRASH view shows a list of all application containers
that have crashed.  A crash normally happens when the application 
running in the container stops or exits unexpectedly.  Often an exit
status code is recorded in the EXIT_DESCRIPTION which may help
understand what happened.  

Java exit status codes:
    137 = OutOfMemory error
    143 = SIGTERM (sometimes as a result of OutOfMemory error)
    255 = OutOfMemory, file descriptors, other error
`

const HelpColumnsText = `
 **CRASH List Columns:**

  CRASH_TIME - Date/time of the crash (24 hour format in local timezone)
  IDX - Application container index that crashed
  EXIT_DESCRIPTION - Container exit desscription often showing exit code  
`

const HelpLocalViewKeybindings = `
`
