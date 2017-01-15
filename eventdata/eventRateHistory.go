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

package eventdata

import (
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/mohae/deepcopy"
)

type EventRateDetail struct {
	RateHigh int
	RateLow  int
	RateAvg  int
}

type EventRate struct {
	BeginTime          time.Time
	EndTime            time.Time
	EventRateDetailMap map[events.Envelope_EventType]*EventRateDetail
	TotalHigh          int
	TotalLow           int
	TotalAvg           int
}

type HistoryRangeType int

const (
	BY_SECOND HistoryRangeType = iota
	BY_MINUTE
	BY_10MINUTE
	BY_HOUR
	BY_DAY
)

// Consolidate history data as follows
// The following configuration results is a max of 764 records in 7 days
// Then just 1 additional record for every day after that.
var recordTimeRangeMaximums = []struct {
	TimeRangeType HistoryRangeType
	RecordMax     int
}{
	{BY_SECOND, 60 * 2},    // Keep 60 seconds minimum of second resolution
	{BY_MINUTE, 10 * 2},    // Keep 60 minutes minimum of minute resolution
	{BY_10MINUTE, 144 * 2}, // Keep 144 10-min records (24 hours) minimum of 10-min resolution
	{BY_HOUR, 168 * 2},     // Keep 168 hours (7 days) minimum of 1 hour resolution
	{BY_DAY, 0},            // Keep 1 day resolution forever
}

type EventRateHistory struct {
	eventProcessor         *EventProcessor
	eventRateByDurationMap map[HistoryRangeType][]*EventRate
	lastTimeHistoryCapture time.Time
	mu                     sync.Mutex

	// This is used when display is "paused" to allow a snapshot
	// of data at the time of the cause
	frozenEventRate []*EventRate
}

func NewEventRateHistory(ep *EventProcessor) *EventRateHistory {
	erh := &EventRateHistory{
		eventProcessor:         ep,
		eventRateByDurationMap: make(map[HistoryRangeType][]*EventRate),
	}
	return erh
}

func (erh *EventRateHistory) GetCurrentHistory() []*EventRate {
	// Merge all the buckets
	// TODO: Might need to thread product this
	erh.mu.Lock()
	defer erh.mu.Unlock()
	mergedHistory := make([]*EventRate, 0)
	for _, eventRate := range erh.eventRateByDurationMap {
		mergedHistory = append(mergedHistory, eventRate...)
	}
	return mergedHistory
}

func (erh *EventRateHistory) GetDisplayedHistory() []*EventRate {
	if erh.frozenEventRate != nil {
		return erh.frozenEventRate
	} else {
		return erh.GetCurrentHistory()
	}
}

// SetFreezeData is called when user pauses/unpauses display
func (erh *EventRateHistory) SetFreezeData(freezeData bool) {
	if freezeData {
		erh.frozenEventRate = deepcopy.Copy(erh.GetCurrentHistory()).([]*EventRate)
	} else {
		erh.frozenEventRate = nil
	}
}

func (erh *EventRateHistory) GetCurrentRate() int {
	erh.mu.Lock()
	defer erh.mu.Unlock()
	currentRate := 0
	eventRateList := erh.eventRateByDurationMap[BY_SECOND]
	if eventRateList != nil && len(eventRateList) > 0 {
		eventRate := eventRateList[len(eventRateList)-1]
		currentRate = eventRate.TotalHigh
	}
	return currentRate
}

func (erh *EventRateHistory) start() {
	ticker := time.NewTicker(time.Second)
	erh.lastTimeHistoryCapture = time.Now()
	go func() {
		toplog.Info("EventRateHistory tracking started")
		for t := range ticker.C {
			erh.captureCurrentRates(t)
		}
	}()
}

func (erh *EventRateHistory) captureCurrentRates(time time.Time) {

	erh.mu.Lock()
	defer erh.mu.Unlock()

	er := &EventRate{
		EventRateDetailMap: make(map[events.Envelope_EventType]*EventRateDetail),
	}
	ep := erh.eventProcessor

	er.BeginTime = erh.lastTimeHistoryCapture
	er.EndTime = time
	erh.lastTimeHistoryCapture = time

	for eventType, rateCounter := range ep.eventRateCounterMap {
		erd := &EventRateDetail{}
		er.EventRateDetailMap[eventType] = erd
		rate := rateCounter.Rate()
		erd.RateHigh = rate
		er.TotalHigh += rate
		erd.RateLow = rate
		er.TotalLow += rate
		erd.RateAvg = rate
		er.TotalAvg += rate
	}

	rateBySecondList := erh.eventRateByDurationMap[BY_SECOND]
	erh.eventRateByDurationMap[BY_SECOND] = append(rateBySecondList, er)

	erh.consolidateHistoryData()
}

func (erh *EventRateHistory) consolidateHistoryData() {
	// Consolidate history data as follows (example)
	// 120 max -> 60 consolidated records at 1 second (60 seconds)
	// 120 max -> 60 consolidated records at 1 minute  (60 minutes)
	// 288 max -> 144 consolidated records at 10 minutes (24 hours)
	// 336 max -> 168 consolidated records at 1 hour (7 days)
	// n records at 1 day (> 7 days)
	for i, recordTimeRangeMax := range recordTimeRangeMaximums {

		nextTimeRangeType := BY_DAY
		if recordTimeRangeMax.RecordMax > 0 {
			nextTimeRangeType = recordTimeRangeMaximums[i+1].TimeRangeType
		}
		if !erh.consolidateHistoryDataByTimeRange(recordTimeRangeMax.TimeRangeType, recordTimeRangeMax.RecordMax, nextTimeRangeType) {
			// If no records were consolidated at one level, then nothing needs to be done at next levels
			break
		}
	}

}

// consolidateHistoryDataByTimeRange if records were consolidated return true
func (erh *EventRateHistory) consolidateHistoryDataByTimeRange(timeRangeType HistoryRangeType, maxRecords int, nextTimeRangeType HistoryRangeType) bool {

	rateBySecondList := erh.eventRateByDurationMap[timeRangeType]
	if maxRecords > 0 && len(rateBySecondList) > maxRecords {
		consolidateQuantity := maxRecords / 2
		olderRecords := rateBySecondList[0:consolidateQuantity]
		newerRecords := rateBySecondList[consolidateQuantity:len(rateBySecondList)]
		erh.eventRateByDurationMap[timeRangeType] = newerRecords

		eventRate := erh.createConsolidatedEventRate(olderRecords)
		rateByMinuteList := erh.eventRateByDurationMap[nextTimeRangeType]
		erh.eventRateByDurationMap[nextTimeRangeType] = append(rateByMinuteList, eventRate)
		return true
	}
	return false
}

func (erh *EventRateHistory) createConsolidatedEventRate(eventRateList []*EventRate) *EventRate {
	consolidatedEventRate := &EventRate{}
	consolidatedEventRate.BeginTime = eventRateList[0].BeginTime
	consolidatedEventRate.EndTime = eventRateList[len(eventRateList)-1].EndTime
	consolidatedEventRate.EventRateDetailMap = make(map[events.Envelope_EventType]*EventRateDetail)

	for _, eventRate := range eventRateList {

		SetMax(&consolidatedEventRate.TotalHigh, eventRate.TotalHigh)
		SetMin(&consolidatedEventRate.TotalLow, eventRate.TotalLow)

		for eventType, eventRateDetail := range eventRate.EventRateDetailMap {
			consolidatedEventRateDetail := consolidatedEventRate.EventRateDetailMap[eventType]
			if consolidatedEventRateDetail == nil {
				consolidatedEventRateDetail = &EventRateDetail{}
				consolidatedEventRate.EventRateDetailMap[eventType] = consolidatedEventRateDetail
			}
			consolidatedEventRate.EventRateDetailMap[eventType] = consolidatedEventRateDetail
			//toplog.Info("createConsolidatedEventRate: %v", eventRateDetail.RateHigh)
			SetMax(&consolidatedEventRateDetail.RateHigh, eventRateDetail.RateHigh)
			SetMin(&consolidatedEventRateDetail.RateLow, eventRateDetail.RateLow)
		}
	}
	return consolidatedEventRate
}

func SetMin(x *int, y int) {
	if *x > y {
		*x = y
	}
}

func SetMax(x *int, y int) {
	if *x < y {
		*x = y
	}
}

func Min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
