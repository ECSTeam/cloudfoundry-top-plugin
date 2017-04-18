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

package eventRateHistoryView

import "github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon/views/helpView"

const HelpText = HelpOverviewText + helpView.HelpHeaderText + HelpColumnsText + helpView.HelpTopLevelDataViewKeybindings + helpView.HelpCommonDataViewKeybindings

const HelpOverviewText = `
**Event Rate Peak History View**

Event rate peak history shows the peak rate of events per second coming
from the firehose. This peak rate is measured from the perspective of
the top client. Two instance of top running against the same foundation
can record different peak events because of network hops and process
scheduling. To conserve memory, history data is consolidated as follows:

  Keep 60 seconds minimum of second resolution
  Keep 60 minutes minimum of minute resolution
  Keep 144 10-min records (24 hours) minimum of 10-min resolution
  Keep 168 hours (7 days) minimum of 1 hour resolution
  Keep 1 day resolution forever

This configuration results is a max of 764 records in 7 days then just
1 additional record for every day after that.

`

const HelpColumnsText = `
**Event Rate Peak Columns:**

  BEGIN_TIME - Date/Time of when the peak rate measurment was started
  END_TIME - Date/Time of when the peak rate measurment was ended
  INTR - The interval, in seconds, of the history record (end-begin time)
  TOTAL - Total events per second that occured during the time interval
  HTTP - Peak HTTP events per second that occured during the time
       interval
  CONTAINER - Peak CONTAINER events per second that occured during the
       time interval
  LOG - Peak LOG events per second that occured during the time interval
  VALUE - Peak VALUE events per second that occured during the time
       interval
  COUNTER - Peak COUNTER events per second that occured during the time
       interval
  ERROR - Peak ERROR events per second that occured during the time
       interval
`
