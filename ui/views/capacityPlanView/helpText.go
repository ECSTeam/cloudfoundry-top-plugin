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

package capacityPlanView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Capacity Plan View**

Capacity plan view shows how many containers of various memory sizes
can be deployed to the foundation based on current capacity.  

NOTE: 
This view currently does not take into consideration the stack
(e.g., cflinuxfs2 vs windows2012R2 ).  This will be a future
enhancement.
`
const HelpColumnsText = `
**Capacity Plan Columns:**

  CELL_IP - IP address of Cloud Foundry cell
  CPUS - Number of CPUs in cell VM
  TOT_MEM - Total Memory in cell VM available for containers
  FREE_MEM - Free Memory in cell VM available for containers
  C_RSVD_MEM - Memory reserved by all containers on cell
  C_USD_MEM - Memory actually in use by all containers
  MAX_CNTR - Max containers a cell can handle
  CNTRS - Number of containers running on cell reported by cell
  0.5GB - Number of 500Meg containers that could be deployed to foundation
  1.0GB - Number of 1GB containers that could be deployed to foundation
  1.5GB - Number of 1.5GB containers that could be deployed to foundation
  2.0GB - Number of 2GB containers that could be deployed to foundation
  2.5GB - Number of 2.5GB containers that could be deployed to foundation
  3.0GB - Number of 3GB containers that could be deployed to foundation
  3.5GB - Number of 3.5GB containers that could be deployed to foundation
  4.0GB - Number of 3GB containers that could be deployed to foundation
  `
