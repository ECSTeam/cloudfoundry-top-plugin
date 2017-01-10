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

package cellView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Cell Stats View**

Cell list view shows a list of all diego cells in the foundation.  This
list may not be complete until the warm-up period is complete.
`

const HelpColumnsText = `
**Cell Columns:**

  CELL_IP - IP address of Cloud Foundry diego cell
  CPU%% - CPU percent consumed by all containers on cell
  RCR - Reporting containers
  CPUS - Number of CPUs in cell VM
  TOT_MEM - Total Memory in cell VM available for containers
  FREE_MEM - Free Memory in cell VM available for containers
  C_RSVD_MEM - Memory reserved by all containers on cell
  C_USD_MEM - Memory actually in use by all containers
  TOT_DISK - Total Disk space in cell VM
  FREE_DISK - Free Disk space in cell VM available for containers
  C_RSVD_DSK - Total disk reserved by all containers on cell
  C_USD_DSK - Total disk actually in use by all containers
  MAX_CNTR - Max containers a cell can handle
  CNTRS - Number of containers running on cell reported by cell
  DNAME - BOSH deployment name
  JOB_NAME - BOSH job name
  JOB_IDX - BOSH job index
`
