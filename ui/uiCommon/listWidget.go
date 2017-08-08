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

package uiCommon

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/ansel1/merry"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
	"github.com/jroimartin/gocui"
)

const (
	// Unicode characters: http://graphemica.com/unicode/characters/page/34
	DownArrow       = string('\U00002193')
	UpArrow         = string('\U00002191')
	DownArrowTiny   = string('\U0000A71C')
	UpArrowTiny     = string('\U0000A71B')
	TriangleUp      = string('\U000025B4')
	TriangleDown    = string('\U000025BE')
	RightArrow      = string('\U00002192')
	LeftArrow       = string('\U00002190')
	InfoIcon        = string('\U00002139')
	Ellipsis        = string('\U00002026')
	TwoDot          = string('\U00002025')
	OneDot          = string('\U00002024')
	CircleBackslash = string('\U000020E0')
)

type preRowDisplayFunc func(data IData, isSelected bool) string
type getRowDisplayFunc func(data IData, columnOwner IColumnOwner) string
type getRowRawValueFunc func(data IData) string
type getDisplayHeaderFunc func() string
type getRowAttentionFunc func(data IData, columnOwner IColumnOwner) AttentionType
type IColumnOwner interface{}

type changeSelectionCallbackFunc func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool

type DisplayViewInterface interface {
	RefreshDisplay(g *gocui.Gui) error
	GetTopOffset() int
}

type ColumnType int

const (
	ALPHANUMERIC ColumnType = iota
	NUMERIC
	TIMESTAMP
)

// Used to determine the attention level of each table's cell (specific field in a display table)
type AttentionType int

const (
	ATTENTION_NORMAL AttentionType = iota
	ATTENTION_HOT
	ATTENTION_WARM
	ATTENTION_NOT_DESIRED_STATE
	ATTENTION_ACTIVITY
	ATTENTION_ALERT
	ATTENTION_WARN
	ATTENTION_NOT_MONITORED
	ATTENTION_STATE_STARTING
	ATTENTION_STATE_UNKNOWN
	ATTENTION_STATE_DOWN
	ATTENTION_STATE_TERM
	ATTENTION_STATE_CRASHED
	ATTENTION_CONTAINER_SHORT_UPTIME
)

type ListColumn struct {
	id                 string
	label              string
	size               int
	columnType         ColumnType
	leftJustifyLabel   bool
	sortFunc           util.LessFunc
	defaultReverseSort bool
	displayFunc        getRowDisplayFunc
	rawValueFunc       getRowRawValueFunc
	attentionFunc      getRowAttentionFunc
}

const LOCK_COLUMNS = 1

type IData interface {
	Id() string
}

type ListWidget struct {
	masterUI     masterUIInterface.MasterUIInterface
	name         string
	bottomMargin int

	Title string

	displayView DisplayViewInterface

	highlightKey          string
	displayRowIndexOffset int
	displayColIndexOffset int

	PreRowDisplayFunc  preRowDisplayFunc
	columnOwner        IColumnOwner
	listData           []IData
	unfilteredListData []IData

	columns   []*ListColumn
	columnMap map[string]*ListColumn

	selectColumnMode bool
	selectedColumnId string

	sortColumns []*SortColumn

	filterColumnMap map[string]*FilterColumn
}

type SortColumn struct {
	Id          string
	ReverseSort bool
}

type FilterColumn struct {
	filterText    string
	compiledRegex *regexp.Regexp
}

var (
	normalHeaderColor    string
	savedSortColumns     map[string][]*SortColumn
	savedFilterColumnMap map[string]map[string]*FilterColumn
)

func init() {
	if util.IsMSWindows() {
		// Windows cmd.exe supports only 8 colors
		normalHeaderColor = util.DIM_WHITE
	} else {
		normalHeaderColor = util.WHITE_TEXT_SOFT_BG
	}
	// savedSortColumns map key = viewName
	savedSortColumns = make(map[string][]*SortColumn)
	// savedFilterColumnMap map key = viewName
	savedFilterColumnMap = make(map[string]map[string]*FilterColumn)
}

func NewSortColumn(id string, reverseSort bool) *SortColumn {
	return &SortColumn{Id: id, ReverseSort: reverseSort}
}

func NewListColumn(
	id, label string,
	size int,
	columnType ColumnType,
	leftJustifyLabel bool,
	sortFunc util.LessFunc,
	defaultReverseSort bool,
	displayFunc getRowDisplayFunc,
	rawValueFunc getRowRawValueFunc,
	attentionFunc getRowAttentionFunc) *ListColumn {
	column := &ListColumn{
		id:                 id,
		label:              label,
		size:               size,
		columnType:         columnType,
		leftJustifyLabel:   leftJustifyLabel,
		sortFunc:           sortFunc,
		defaultReverseSort: defaultReverseSort,
		displayFunc:        displayFunc,
		rawValueFunc:       rawValueFunc,
		attentionFunc:      attentionFunc,
	}

	return column
}

func NewListWidget(masterUI masterUIInterface.MasterUIInterface, name string,
	bottomMargin int, displayView DisplayViewInterface,
	columns []*ListColumn, columnOwner IColumnOwner, defaultSortColumns []*SortColumn) *ListWidget {
	w := &ListWidget{
		masterUI:        masterUI,
		name:            name,
		bottomMargin:    bottomMargin,
		displayView:     displayView,
		columns:         columns,
		columnMap:       make(map[string]*ListColumn),
		filterColumnMap: make(map[string]*FilterColumn),
		columnOwner:     columnOwner,
	}
	for _, col := range columns {
		w.columnMap[col.id] = col
	}

	sortColumns := savedSortColumns[name]
	if sortColumns == nil {
		sortColumns = defaultSortColumns
	}
	w.SetSortColumns(sortColumns)

	savedFilterColumnMap := savedFilterColumnMap[name]
	if savedFilterColumnMap != nil {
		w.filterColumnMap = savedFilterColumnMap
	}

	return w
}

func (w *ListWidget) Name() string {
	return w.name
}

func (w *ListWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	bottom := maxY - w.bottomMargin
	topMargin := w.displayView.GetTopOffset()
	if topMargin >= bottom {
		bottom = topMargin + 1
	}

	// Check if the view has been resized, if not then do nothing
	// This prevents non-visable views from doing more work then needed
	v, err := g.View(w.name)
	if err == nil {
		x, y := v.Size()
		if maxX-2 == x && bottom-topMargin-1 == y {
			//toplog.Info("x:%v  maxX-1: %v  y: %v  bottom-topMargin: %v  ", x, maxX-2, y, bottom-topMargin-1)
			return nil
		}
	}

	v, err = g.SetView(w.name, 0, topMargin, maxX-1, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " (ListWidget) layout error:" + err.Error())
		}
		v.Title = w.Title
		v.Frame = true
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgdn, gocui.ModNone, w.pageDownAction); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgup, gocui.ModNone, w.pageUpAction); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.arrowRight); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.arrowLeft); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyHome, gocui.ModNone, w.arrowHome); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnd, gocui.ModNone, w.arrowEnd); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(w.name, 'o', gocui.ModNone, w.editSortAction); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(w.name, 'f', gocui.ModNone, w.editFilterAction); err != nil {
			log.Panicln(err)
		}

		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				w.highlightKey = ""
				w.displayRowIndexOffset = 0
				w.RefreshDisplay(g)
				return nil
			}); err != nil {
			log.Panicln(err)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	return w.RefreshDisplay(g)
}

func (asUI *ListWidget) HighlightKey() string {
	return asUI.highlightKey
}

// Get the highlighted data row
func (asUI *ListWidget) HighlightData() IData {
	for _, data := range asUI.unfilteredListData {
		if data.Id() == asUI.highlightKey {
			return data
		}
	}
	return nil
}

func (asUI *ListWidget) GetFilterColumnMap() map[string]*FilterColumn {
	return asUI.filterColumnMap
}

func (asUI *ListWidget) SetListData(listData []IData) {
	asUI.unfilteredListData = listData
	asUI.FilterAndSortData()
}

func (asUI *ListWidget) GetListData() []IData {
	return asUI.unfilteredListData
}

func (asUI *ListWidget) FilterAndSortData() {
	filteredData := asUI.filterData(asUI.unfilteredListData)
	asUI.listData = asUI.sortData(filteredData)
}

func (asUI *ListWidget) sortData(listData []IData) []IData {
	sortFunctions := asUI.GetSortFunctions()
	sortData := make([]util.Sortable, 0, len(listData))
	//toplog.Debug("sortStats size before:%v", len(sortStats))
	for _, data := range listData {
		sortData = append(sortData, data)
	}
	//toplog.Debug("sortStats size after:%v", len(sortStats))
	util.OrderedBy(sortFunctions).Sort(sortData)

	s2 := make([]IData, len(sortData))
	for i, d := range sortData {
		s2[i] = d.(IData)
	}
	return s2
}

func (asUI *ListWidget) filterData(listData []IData) []IData {
	filteredList := make([]IData, 0, len(listData))
	for _, data := range listData {
		if asUI.FilterRow(data) {
			filteredList = append(filteredList, data)
		}
	}
	return filteredList
}

func (asUI *ListWidget) FilterRow(data IData) bool {
	for _, column := range asUI.columns {
		filter := asUI.filterColumnMap[column.id]
		if filter != nil && filter.filterText != "" {
			if !asUI.filterRow(data, column, filter) {
				return false
			}
		}
	}
	return true
}

func (asUI *ListWidget) filterRow(data IData, column *ListColumn, filter *FilterColumn) bool {
	switch column.columnType {
	case ALPHANUMERIC:
		return asUI.filterRowAlpha(data, column, filter)
	case NUMERIC:
		return asUI.filterRowNumeric(data, column, filter)
	}
	return false
}

func (asUI *ListWidget) filterRowNumeric(data IData, column *ListColumn, filter *FilterColumn) bool {
	varName := "VALUE"
	expressionStr := filter.filterText
	if expressionStr == "" {
		return true
	}
	expressionStr = varName + " " + expressionStr
	value := column.rawValueFunc(data)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return true
	}

	expression, err := govaluate.NewEvaluableExpression(expressionStr)
	if err != nil {
		return true
	}

	parameters := make(map[string]interface{}, 8)
	parameters[varName] = floatValue

	result, err := expression.Evaluate(parameters)
	if err != nil {
		return true
	}
	return result.(bool)
}

func (asUI *ListWidget) filterRowAlpha(data IData, column *ListColumn, filter *FilterColumn) bool {
	regex := filter.compiledRegex
	if regex == nil {
		// make regex case insenstive
		filterText := "(?i)" + filter.filterText
		compiledRegex, err := regexp.Compile(filterText)
		if err != nil {
			// Better to error on the side of showing the row
			return true
		}
		filter.compiledRegex = compiledRegex
		regex = compiledRegex
	}
	value := column.rawValueFunc(data)

	return regex.MatchString(value)
}

func (asUI *ListWidget) SaveFilters() {
	savedFilterColumnMap[asUI.name] = asUI.filterColumnMap
}

func (asUI *ListWidget) SetSortColumns(sortColumns []*SortColumn) {
	asUI.sortColumns = sortColumns
	// TOOD: The savedSortColumns should write to json file to save the sort order
	// between runs of top plugin
	savedSortColumns[asUI.name] = sortColumns
}

func (asUI *ListWidget) GetSortColumns() []*SortColumn {
	return asUI.sortColumns
}

func (asUI *ListWidget) GetSortFunctions() []util.LessFunc {

	sortFunctions := make([]util.LessFunc, 0)
	for _, sortColumn := range asUI.sortColumns {
		sc := asUI.columnMap[sortColumn.Id]
		if sc == nil {
			log.Panic(merry.Errorf("Unable to find sort column: %v", sortColumn.Id))
		}
		sortFunc := sc.sortFunc
		if sortColumn.ReverseSort {
			sortFunc = util.Reverse(sortFunc)
		}
		sortFunctions = append(sortFunctions, sortFunc)
	}
	return sortFunctions
}

func (asUI *ListWidget) GetColumns() []*ListColumn {
	return asUI.columns
}

func (asUI *ListWidget) RefreshDisplay(g *gocui.Gui) error {

	v, err := g.View(asUI.name)
	if err != nil {
		return err
	}
	_, maxY := v.Size()
	maxRows := maxY - 1

	title := asUI.Title
	displayListSize := len(asUI.listData)
	unfilteredListSize := len(asUI.unfilteredListData)
	if displayListSize != unfilteredListSize {
		title = fmt.Sprintf("%v (filter showing %v of %v)", title, displayListSize, unfilteredListSize)
	}
	v.Title = title

	v.Clear()
	listSize := len(asUI.listData)
	offset := asUI.displayRowIndexOffset
	if offset > listSize || offset < 0 {
		// If the list changes size, specifically gets smaller since the last time we've
		// refreshed AND we're scroll beyond the new list size, reset display offset back
		// to zero
		asUI.displayRowIndexOffset = 0
		offset = 0
	}

	if listSize > 0 || asUI.selectColumnMode {
		stopRowIndex := maxRows + offset
		asUI.writeHeader(g, v)

		/*
			toplog.Info("listWidget listSize: %v  stopRowIndex: %v  maxRows: %v  offset: %v",
				listSize, stopRowIndex, maxRows, offset)
		*/

		// Loop through all rows
		for i := 0; i < listSize && i < stopRowIndex; i++ {
			if i < offset {
				continue
			}
			asUI.writeRowData(g, v, i)
		}
	} else {
		if len(asUI.unfilteredListData) > 0 {
			fmt.Fprint(v, " \n No data to display because of filters")
		} else {
			fmt.Fprint(v, " \n No data yet...")
		}
	}

	return nil
}

func (asUI *ListWidget) writeRowData(g *gocui.Gui, v *gocui.View, rowIndex int) {
	rowData := asUI.listData[rowIndex]
	isSelected := false
	if rowData.Id() == asUI.highlightKey {
		fmt.Fprint(v, util.REVERSE_GREEN)
		isSelected = true
	}

	sortColumnId := ""
	if len(asUI.sortColumns) > 0 {
		sortColumnId = asUI.sortColumns[0].Id
	}

	if asUI.PreRowDisplayFunc != nil {
		fmt.Fprint(v, asUI.PreRowDisplayFunc(rowData, isSelected))
	}

	// Loop through all columns
	for colIndex, column := range asUI.columns {
		colorString := ""
		if colIndex > asUI.lastColumnCanDisplay(g, asUI.displayColIndexOffset) {
			break
		}
		if colIndex >= LOCK_COLUMNS && colIndex < asUI.displayColIndexOffset+LOCK_COLUMNS {
			continue
		}

		if !isSelected && column.attentionFunc != nil {
			attentionLevel := column.attentionFunc(rowData, asUI.columnOwner)
			attributeModifier := ""
			if column.id == sortColumnId {
				attributeModifier = util.BRIGHT
			} else {
				attributeModifier = util.DIM
			}
			switch attentionLevel {
			case ATTENTION_HOT:
				colorString = util.RED
			case ATTENTION_WARM:
				colorString = util.YELLOW
			case ATTENTION_NOT_DESIRED_STATE:
				colorString = util.RED
			case ATTENTION_ALERT:
				colorString = util.RED
			case ATTENTION_WARN:
				colorString = util.YELLOW
			case ATTENTION_ACTIVITY:
				colorString = util.CYAN
			case ATTENTION_NOT_MONITORED:
				colorString = util.BRIGHT_BLACK
			case ATTENTION_STATE_STARTING:
				colorString = util.YELLOW
			case ATTENTION_STATE_UNKNOWN:
				colorString = util.YELLOW
			case ATTENTION_STATE_DOWN:
				colorString = util.RED
			case ATTENTION_STATE_TERM:
				colorString = util.PURPLE
			case ATTENTION_STATE_CRASHED:
				colorString = util.RED
			case ATTENTION_CONTAINER_SHORT_UPTIME:
				colorString = util.CYAN
			}
			if len(colorString) == 4 {
				colorString = colorString + attributeModifier
			}
		}

		if !isSelected && colorString == "" && column.id == sortColumnId {
			colorString = util.BRIGHT_WHITE
		}
		if colorString != "" {
			fmt.Fprintf(v, "%v", colorString)
		}

		fmt.Fprint(v, column.displayFunc(rowData, asUI.columnOwner))
		if !isSelected && colorString != "" {
			fmt.Fprint(v, util.CLEAR)
		}
		fmt.Fprint(v, " ")
	}
	fmt.Fprint(v, "\n")
	fmt.Fprint(v, util.CLEAR)
}

func (asUI *ListWidget) writeHeader(g *gocui.Gui, v *gocui.View) {

	lastColumnCanDisplay := asUI.lastColumnCanDisplay(g, asUI.displayColIndexOffset)

	fmt.Fprint(v, normalHeaderColor)

	// Loop through all columns (for headers)
	for colIndex, column := range asUI.columns {
		if colIndex > lastColumnCanDisplay {
			break
		}
		if colIndex >= LOCK_COLUMNS && colIndex < asUI.displayColIndexOffset+LOCK_COLUMNS {
			continue
		}
		colorString := ""
		editSortColumn := false
		if asUI.selectColumnMode && asUI.selectedColumnId == column.id {
			editSortColumn = true
			fmt.Fprint(v, util.REVERSE_WHITE)
		}
		var buffer bytes.Buffer
		buffer.WriteString("%")
		if column.leftJustifyLabel {
			buffer.WriteString("-")
		}
		buffer.WriteString(strconv.Itoa(column.size))
		buffer.WriteString("v ")

		label := column.label

		if len(asUI.sortColumns) > 0 {
			sortCol := asUI.sortColumns[0]
			if sortCol != nil && sortCol.Id == column.id {
				if !editSortColumn {
					colorString = util.BRIGHT_WHITE
				}
				if sortCol.ReverseSort {
					label = label + DownArrow
				} else {
					label = label + UpArrow
				}
			}
		}

		if colorString != "" {
			fmt.Fprint(v, colorString)
		}
		fmt.Fprintf(v, buffer.String(), label)
		if editSortColumn || colorString != "" {
			fmt.Fprint(v, normalHeaderColor)
		}

	}
	fmt.Fprint(v, util.CLEAR)
	fmt.Fprint(v, "\n")
}

func (asUI *ListWidget) lastColumnCanDisplay(g *gocui.Gui, ifDisplayColIndexOffset int) int {

	v, err := g.View(asUI.name)
	if err != nil {
		return 0
	}

	maxX, _ := v.Size()
	totalWidth := 0
	lastColumnCanDisplay := len(asUI.columns) - 1
	for colIndex, column := range asUI.columns {
		if colIndex >= LOCK_COLUMNS && colIndex < ifDisplayColIndexOffset+LOCK_COLUMNS {
			continue
		}
		totalWidth = totalWidth + column.size
		if totalWidth > maxX {
			lastColumnCanDisplay = colIndex - 1
			break
		}
		// Add one for the space after the column name
		totalWidth = totalWidth + 1
	}
	if lastColumnCanDisplay < 0 {
		// Must display at least 1 column
		lastColumnCanDisplay = 0
	}
	//writeFooter(g, fmt.Sprintf("\r ** maxX:%v lastColumnCanDisplay:%v totalWidth:%v", maxX,lastColumnCanDisplay,totalWidth))
	return lastColumnCanDisplay
}

func (asUI *ListWidget) scollSelectedColumnIntoView(g *gocui.Gui) error {
	offsetFromView, indexOfSelectedCol := asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)

	if offsetFromView < 0 {
		asUI.displayColIndexOffset = asUI.displayColIndexOffset + offsetFromView
	} else {

		newDisplayOffset := asUI.displayColIndexOffset + offsetFromView
		for asUI.lastColumnCanDisplay(g, newDisplayOffset) < indexOfSelectedCol {
			newDisplayOffset++
		}
		asUI.displayColIndexOffset = newDisplayOffset
	}
	return asUI.RefreshDisplay(g)

}

func (asUI *ListWidget) isColumnVisable(g *gocui.Gui, findColumnId string) bool {
	offsetFromView, _ := asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)
	return offsetFromView == 0
}

// columnOffsetFromVisability returns the number of columns to the left or right
// of current display.  E.g., -2 means its 2 index positions to the left.
// returns:
// columnOffsetFromVisability
// index position of findColumnId
func (asUI *ListWidget) columnOffsetFromVisability(g *gocui.Gui, findColumnId string) (int, int) {
	lastColumnCanDisplay := asUI.lastColumnCanDisplay(g, asUI.displayColIndexOffset)
	firstDisplayedColumnIndex := -1
	foundColumnIndex := 0
	// Loop through all columns (for headers)
	for colIndex, column := range asUI.columns {
		if colIndex < asUI.displayColIndexOffset+LOCK_COLUMNS {
			//continue
		} else {
			if firstDisplayedColumnIndex == -1 {
				firstDisplayedColumnIndex = colIndex
			}
		}
		if findColumnId == column.id {
			foundColumnIndex = colIndex
		}
	}
	if foundColumnIndex < firstDisplayedColumnIndex {
		// Return negative offset from visability on left
		return foundColumnIndex - firstDisplayedColumnIndex, foundColumnIndex
	} else if foundColumnIndex > lastColumnCanDisplay {
		// Return positive offset from visability on right
		return foundColumnIndex - lastColumnCanDisplay, foundColumnIndex
	}
	return 0, foundColumnIndex
}

func (asUI *ListWidget) arrowRight(g *gocui.Gui, v *gocui.View) error {
	lastColumnCanDisplay := asUI.lastColumnCanDisplay(g, asUI.displayColIndexOffset)
	if lastColumnCanDisplay < len(asUI.columns)-1 {
		asUI.displayColIndexOffset++
	}
	asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowLeft(g *gocui.Gui, v *gocui.View) error {
	asUI.displayColIndexOffset--
	if asUI.displayColIndexOffset < 0 {
		asUI.displayColIndexOffset = 0
	}
	asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowHome(g *gocui.Gui, v *gocui.View) error {
	asUI.displayColIndexOffset = 0
	asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowEnd(g *gocui.Gui, v *gocui.View) error {
	for {
		lastColumnCanDisplay := asUI.lastColumnCanDisplay(g, asUI.displayColIndexOffset)
		if lastColumnCanDisplay < len(asUI.columns)-1 {
			asUI.displayColIndexOffset++
		} else {
			break
		}
	}

	asUI.columnOffsetFromVisability(g, asUI.selectedColumnId)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowUp(g *gocui.Gui, v *gocui.View) error {
	listSize := len(asUI.listData)
	callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
		if rowIndex > 0 {
			_, viewY := v.Size()
			viewSize := viewY - 1
			asUI.highlightKey = lastKey
			offset := rowIndex - 1
			if listSize > viewSize && offset > listSize-viewSize {
				offset = listSize - viewSize
			}
			if asUI.displayRowIndexOffset > offset || rowIndex > asUI.displayRowIndexOffset+viewSize {
				asUI.displayRowIndexOffset = offset
			}
			return true
		}
		return false
	}
	return asUI.moveHighlight(g, v, callbackFunc)
}

func (asUI *ListWidget) arrowDown(g *gocui.Gui, v *gocui.View) error {

	listSize := len(asUI.listData)
	callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
		if rowIndex+1 < listSize {
			_, viewY := v.Size()
			offset := (rowIndex + 2) - (viewY - 1)
			if offset > asUI.displayRowIndexOffset || rowIndex < asUI.displayRowIndexOffset {
				asUI.displayRowIndexOffset = offset
			}
			asUI.highlightKey = asUI.listData[rowIndex+1].Id()
			return true
		}
		return false
	}
	return asUI.moveHighlight(g, v, callbackFunc)
}

func (asUI *ListWidget) moveHighlight(g *gocui.Gui, v *gocui.View, callback changeSelectionCallbackFunc) error {

	listSize := len(asUI.listData)
	if asUI.highlightKey == "" {
		if listSize > 0 {
			asUI.highlightKey = asUI.listData[0].Id()
		}
	} else {
		lastKey := ""
		foundMatch := false
		for rowIndex := 0; rowIndex < listSize; rowIndex++ {
			if asUI.listData[rowIndex].Id() == asUI.highlightKey {
				foundMatch = true
				if callback(g, v, rowIndex, lastKey) {
					break
				}
			}
			lastKey = asUI.listData[rowIndex].Id()
		}
		if !foundMatch {
			if listSize > 0 {
				asUI.highlightKey = asUI.listData[0].Id()
			}
		}

	}
	return asUI.RefreshDisplay(g)

}

func (asUI *ListWidget) pageUpAction(g *gocui.Gui, v *gocui.View) error {
	callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
		if rowIndex > 0 {
			_, viewY := v.Size()
			viewSize := viewY - 1
			offset := 0
			if rowIndex == asUI.displayRowIndexOffset {
				offset = rowIndex - viewSize
			} else {
				offset = rowIndex - (rowIndex - asUI.displayRowIndexOffset)
			}
			if offset < 0 {
				offset = 0
			}
			if offset > len(asUI.listData) {
				offset = len(asUI.listData) - 1
			}
			if offset < asUI.displayRowIndexOffset {
				asUI.displayRowIndexOffset = offset
			}
			asUI.highlightKey = asUI.listData[offset].Id()
			return true
		}
		return false
	}
	return asUI.moveHighlight(g, v, callbackFunc)
}

func (asUI *ListWidget) pageDownAction(g *gocui.Gui, v *gocui.View) error {
	listSize := len(asUI.listData)
	callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
		if rowIndex < listSize {
			_, viewY := v.Size()
			viewSize := viewY - 1
			offset := 0
			if rowIndex == (asUI.displayRowIndexOffset + viewSize - 1) {
				offset = rowIndex + viewSize
			} else {
				offset = rowIndex + (viewSize - (rowIndex - asUI.displayRowIndexOffset)) - 1
			}
			if offset > listSize-1 {
				offset = listSize - 1
			}
			if offset > (asUI.displayRowIndexOffset + viewSize - 1) {
				asUI.displayRowIndexOffset = offset - viewSize + 1
			}
			asUI.highlightKey = asUI.listData[offset].Id()
			return true
		}
		return false
	}
	return asUI.moveHighlight(g, v, callbackFunc)
}

// This is for debugging -- remove it later
func writeFooter(g *gocui.Gui, msg string) {
	v, _ := g.View("footerView")
	fmt.Fprint(v, msg)

}

func (asUI *ListWidget) editFilterAction(g *gocui.Gui, v *gocui.View) error {
	editViewName := asUI.name + ".editFilterView"
	asUI.selectColumnMode = true
	if asUI.selectedColumnId == "" {
		asUI.selectedColumnId = asUI.columns[0].id
	}
	filterView := NewEditFilterView(asUI.masterUI, editViewName, asUI)
	asUI.masterUI.LayoutManager().Add(filterView)
	asUI.masterUI.SetCurrentViewOnTop(g)

	// TODO: Is this the correct spot to do this?
	asUI.masterUI.SetEditColumnMode(g, true)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) editSortAction(g *gocui.Gui, v *gocui.View) error {
	editViewName := asUI.name + ".editSortView"
	asUI.selectColumnMode = true
	if asUI.selectedColumnId == "" {
		asUI.selectedColumnId = asUI.columns[0].id
	}
	editView := NewEditSortView(asUI.masterUI, editViewName, asUI)
	asUI.masterUI.LayoutManager().Add(editView)
	asUI.masterUI.SetCurrentViewOnTop(g)
	// TODO: Is this the correct spot to do this?
	asUI.masterUI.SetEditColumnMode(g, true)
	return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) enableSelectColumnMode(enable bool) {
	asUI.selectColumnMode = enable
}
