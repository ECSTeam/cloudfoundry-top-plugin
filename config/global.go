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

package config

// TODO: Future home for storing configuration information to a json file

const MaxTopInternalLogLineHistory = 2000

const WarmUpSeconds = 60
const StaleContainerSeconds = 65
const DeadContainerSeconds = 185

// When an app is started it can take 15+ seconds for the container to report in (RCR)
// In order to not trigger a "application not in desired state" warning too early we
// wait some amount of time from the point where metadata cache shows the app in state "STARTED"
const AppNotInDesiredStateWaitTimeSeconds = 65

// Monitor app details after vist for 15 minutes (900 seconds)
const MonitorAppDetailTTL = 900

const MaxDomainBucket = 100
const MaxHostBucket = 10000
const MaxUserAgentBucket = 100
const MaxForwarderBucket = 100

// Number of records to retrieve per cloud controller REST call.
const ResultsPerPage = 100
const ResultsPerV3Page = 1000
