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

import (
	"fmt"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/uiCommon"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
)

type GetValueFunction func(stats *DisplayEventRateHistoryStats) int

type RateValueType int

const (
	RATE_HIGH RateValueType = iota
	RATE_LOW
	RATE_AVG
)

func columnEventBeginTime() *uiCommon.ListColumn {
	defaultColSize := 22
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayEventRateHistoryStats).BeginTime.Before(c2.(*DisplayEventRateHistoryStats).BeginTime)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayEventRateHistoryStats)
		return fmt.Sprintf("%-22v", stats.BeginTime.Format("01-02-2006 15:04:05"))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayEventRateHistoryStats)
		return fmt.Sprintf("%v", stats.BeginTime)
	}
	c := uiCommon.NewListColumn("BEGIN_TIME", "BEGIN_TIME", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnEventEndTime() *uiCommon.ListColumn {
	defaultColSize := 9
	sortFunc := func(c1, c2 util.Sortable) bool {
		return c1.(*DisplayEventRateHistoryStats).EndTime.Before(c2.(*DisplayEventRateHistoryStats).EndTime)
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayEventRateHistoryStats)
		return fmt.Sprintf("%-9v", stats.EndTime.Format("15:04:05"))
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayEventRateHistoryStats)
		return fmt.Sprintf("%v", stats.EndTime)
	}
	c := uiCommon.NewListColumn("END_TIME", "END_TIME", defaultColSize,
		uiCommon.TIMESTAMP, true, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}

func columnInterval() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int { return int(stats.Duration.Seconds() + 0.25) }
	return columnTemplate("INTR", getValueFunc, 5)
}

func columnTotalRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int { return stats.TotalHigh }
	return columnRateTemplate("TOTAL", getValueFunc)
}

/*
func columnTotalRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int { return stats.TotalLow }
	return columnRateTemplate("TOTAL_LOW", getValueFunc)
}
*/
func columnHttpStartStopEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_HttpStartStop, RATE_HIGH)
	}
	return columnRateTemplate("HTTP", getValueFunc)
}

func columnHttpStartStopEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_HttpStartStop, RATE_LOW)
	}
	return columnRateTemplate("HTTP_LOW", getValueFunc)
}

func columnContainerMetricEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_ContainerMetric, RATE_HIGH)
	}
	return columnRateTemplate("CONTAINER", getValueFunc)
}

func columnContainerMetricEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_ContainerMetric, RATE_LOW)
	}
	return columnRateTemplate("CONTNR_LOW", getValueFunc)
}

func columnLogMessageEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_LogMessage, RATE_HIGH)
	}
	return columnRateTemplate("LOG", getValueFunc)
}
func columnLogMessageEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_LogMessage, RATE_LOW)
	}
	return columnRateTemplate("LOG_LOW", getValueFunc)
}

func columnValueMetricEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_ValueMetric, RATE_HIGH)
	}
	return columnRateTemplate("VALUE", getValueFunc)
}

func columnValueMetricEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_ValueMetric, RATE_LOW)
	}
	return columnRateTemplate("VALUE_LOW", getValueFunc)
}

func columnCounterEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_CounterEvent, RATE_HIGH)
	}
	return columnRateTemplate("COUNTER", getValueFunc)
}
func columnCounterEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_CounterEvent, RATE_LOW)
	}
	return columnRateTemplate("COUNTER_LOW", getValueFunc)
}

func columnErrorEventRateHigh() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_Error, RATE_HIGH)
	}
	return columnRateTemplate("ERROR", getValueFunc)
}
func columnErrorEventRateLow() *uiCommon.ListColumn {
	getValueFunc := func(stats *DisplayEventRateHistoryStats) int {
		return getRateValue(stats, events.Envelope_Error, RATE_LOW)
	}
	return columnRateTemplate("ERROR_LOW", getValueFunc)
}

func getRateValue(stats *DisplayEventRateHistoryStats, eventType events.Envelope_EventType, valueType RateValueType) int {
	rateDetail := stats.EventRateDetailMap[eventType]
	value := 0
	if rateDetail != nil {
		switch valueType {
		case RATE_HIGH:
			value = rateDetail.RateHigh
		case RATE_LOW:
			//value = rateDetail.RateLow
		case RATE_AVG:
			//value = rateDetail.RateAvg
		}
	}
	return value
}

func columnRateTemplate(columnName string, getValueFunc GetValueFunction) *uiCommon.ListColumn {
	return columnTemplate(columnName, getValueFunc, 10)
}

func columnTemplate(columnName string, getValueFunc GetValueFunction, width int) *uiCommon.ListColumn {
	defaultColSize := 12
	sortFunc := func(c1, c2 util.Sortable) bool {
		return getValueFunc(c1.(*DisplayEventRateHistoryStats)) < getValueFunc(c2.(*DisplayEventRateHistoryStats))
	}
	displayFunc := func(data uiCommon.IData, columnOwner uiCommon.IColumnOwner) string {
		stats := data.(*DisplayEventRateHistoryStats)
		value := getValueFunc(stats)
		return util.FormatDisplayDataRight(util.Format(int64(value)), defaultColSize)
	}
	rawValueFunc := func(data uiCommon.IData) string {
		stats := data.(*DisplayEventRateHistoryStats)
		value := getValueFunc(stats)
		return fmt.Sprintf("%v", value)
	}
	c := uiCommon.NewListColumn(columnName, columnName, defaultColSize,
		uiCommon.NUMERIC, false, sortFunc, true, displayFunc, rawValueFunc, nil)
	return c
}
