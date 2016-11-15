package uiCommon

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/Knetic/govaluate"
	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

type EditFilterView struct {
	*EditColumnViewAbs

	labelWidget masterUIInterface.Manager
	inputWidget masterUIInterface.Manager

	oldFilterColumnMap map[string]*FilterColumn

	editField bool
}

func NewEditFilterView(masterUI masterUIInterface.MasterUIInterface, name string, listWidget *ListWidget) *EditFilterView {
	w := &EditFilterView{EditColumnViewAbs: NewEditColumnViewAbs(masterUI, name, listWidget)}
	w.width = 55
	w.height = 14
	w.title = "Edit Filter"

	w.refreshDisplayCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.refreshDisplayCallback(g, v)
	}

	w.initialLayoutCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.initialLayoutCallback(g, v)
	}

	w.applyActionCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.applyActionCallback(g, v)
	}

	w.cancelActionCallbackFunc = func(g *gocui.Gui, v *gocui.View) error {
		return w.cancelActionCallback(g, v)
	}

	// Save old filter for cancel
	w.oldFilterColumnMap = make(map[string]*FilterColumn)
	for columnId, filter := range listWidget.filterColumnMap {
		cloneFilter := &FilterColumn{filterText: filter.filterText}
		w.oldFilterColumnMap[columnId] = cloneFilter
	}

	return w
}

func (w *EditFilterView) initialLayoutCallback(g *gocui.Gui, v *gocui.View) error {

	v.Wrap = true
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keySpaceAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, 'c', gocui.ModNone, w.clearFilterAction); err != nil {
		return err
	}
	return nil
}

func (w *EditFilterView) refreshDisplayCallback(g *gocui.Gui, v *gocui.View) error {

	v.Clear()
	fmt.Fprintln(v, " ")
	selectedColId := w.listWidget.selectedColumnId
	col := w.listWidget.columnMap[selectedColId]
	filter := w.listWidget.filterColumnMap[selectedColId]
	filterText := "--none--"
	if filter != nil && filter.filterText != "" {
		filterText = filter.filterText
	}

	fmt.Fprintf(v, " Column name: %v\n", col.label)
	fmt.Fprintf(v, " Filter: %v\n\n", filterText)

	if w.editField {
		switch col.columnType {
		case ALPHANUMERIC:
			fmt.Fprintf(v, "\n\n RegEx examples:\n")
			fmt.Fprintf(v, " AppA or AppB: appa|appb\n")
			fmt.Fprintf(v, " Starts with 'p' end with 'ch': p([a-z]*)ch\n")
		case NUMERIC:
			fmt.Fprintf(v, "\n\n Expression examples:\n")
			fmt.Fprintf(v, " Greater then: >0.15\n")
			fmt.Fprintf(v, " Equals: ==2   Note the double equals\n")
		}
	} else {
		fmt.Fprintln(v, " RIGHT or LEFT arrow - highlight column")
		fmt.Fprintln(v, " SPACE - select column to edit filter")
		fmt.Fprintln(v, " ENTER - apply filter, ESC to cancel")
		fmt.Fprintln(v, " 'c' - clear all filters")
	}

	return nil
}

func (w *EditFilterView) getSelectedColumn() *ListColumn {
	selectedColId := w.listWidget.selectedColumnId
	col := w.listWidget.columnMap[selectedColId]
	return col
}

func (w *EditFilterView) applyValueCallback(g *gocui.Gui, v *gocui.View, mgr masterUIInterface.Manager, inputValue string) error {

	col := w.getSelectedColumn()
	var err error
	switch col.columnType {
	case ALPHANUMERIC:
		err = w.applyAlphaFilter(g, v, mgr, inputValue)
	case NUMERIC:
		inputValue = w.adjustExpression(inputValue)
		err = w.applyNumericFilter(g, v, mgr, inputValue)
	}

	if err != nil {
		parentView, err2 := g.View(w.name)
		if err2 != nil {
			return err2
		}
		fmt.Fprintf(parentView, "%v", util.BRIGHT_RED)
		fmt.Fprintf(parentView, "\r Error: %v", err)
		fmt.Fprintf(parentView, "%v", util.CLEAR)
		return nil
	}
	selectedColId := w.listWidget.selectedColumnId
	filter := &FilterColumn{filterText: inputValue}
	w.listWidget.filterColumnMap[selectedColId] = filter

	g.Cursor = false
	if err := w.masterUI.CloseView(w.labelWidget); err != nil {
		return err
	}
	if err := w.masterUI.CloseView(w.inputWidget); err != nil {
		return err
	}
	w.editField = false
	w.listWidget.FilterAndSortData()
	w.RefreshDisplay(g)
	return nil
}

func (w *EditFilterView) applyNumericFilter(g *gocui.Gui, v *gocui.View, mgr masterUIInterface.Manager, inputValue string) error {

	varName := "VALUE"
	if inputValue == "" {
		return nil
	}
	inputValue = varName + " " + inputValue

	expression, err := govaluate.NewEvaluableExpression(inputValue)
	if err != nil {
		return err
	}

	parameters := make(map[string]interface{}, 8)
	parameters[varName] = 1.0

	result, err := expression.Evaluate(parameters)
	if err != nil {
		return err
	}
	if reflect.TypeOf(result) != reflect.TypeOf(true) {
		err := errors.New("Expression does not result in boolean result")
		return err
	}
	return nil
}

// Prefix floats with 0 if not digit before decimal  (e.g, .1 becomes 0.1)
func (w *EditFilterView) adjustExpression(value string) string {
	const regex = `[^\d]\.[\d]+`
	//value := strings.Replace(value, "=", "==", -1)
	r := regexp.MustCompile(regex)
	index := r.FindStringSubmatchIndex(value)
	if len(index) > 0 {
		firstPart := value[0 : index[0]+1]
		lastPart := value[index[0]+1 : len(value)]
		value = firstPart + "0" + lastPart
	}
	return value
}

func (w *EditFilterView) applyAlphaFilter(g *gocui.Gui, v *gocui.View, mgr masterUIInterface.Manager, inputValue string) error {
	_, err := regexp.Compile(inputValue)
	if err != nil {
		return err
	}
	return nil
}

func (w *EditFilterView) clearFilterAction(g *gocui.Gui, v *gocui.View) error {
	w.listWidget.filterColumnMap = make(map[string]*FilterColumn)
	w.listWidget.FilterAndSortData()
	w.RefreshDisplay(g)
	return nil
}

func (w *EditFilterView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

	selectedColId := w.listWidget.selectedColumnId
	filter := w.listWidget.filterColumnMap[selectedColId]
	filterText := ""
	if filter != nil {
		filterText = filter.filterText
	}

	labelText := "Filter:"
	maxLength := 30
	valueText := filterText
	topMargin := 4

	w.labelWidget = NewLabel(w, "label", 1, topMargin, labelText)
	cancelCallbackFunc := func(g *gocui.Gui, v *gocui.View) error {
		//fmt.Printf("\n**** CANCELED ****\n")
		//return w.CloseWidget(g, v)
		return nil
	}

	w.inputWidget = NewInput(w, "input", len(labelText)+2, topMargin, maxLength+2,
		maxLength, valueText,
		w.applyValueCallback,
		cancelCallbackFunc)

	w.masterUI.LayoutManager().Add(w)
	w.masterUI.LayoutManager().Add(w.labelWidget)
	w.masterUI.LayoutManager().Add(w.inputWidget)

	w.Layout(g)
	w.labelWidget.Layout(g)
	w.inputWidget.Layout(g)
	w.masterUI.SetCurrentViewOnTop(g, "input")
	w.editField = true

	w.RefreshDisplay(g)

	return nil
}

func (w *EditFilterView) applyActionCallback(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func (w *EditFilterView) cancelActionCallback(g *gocui.Gui, v *gocui.View) error {
	w.listWidget.filterColumnMap = w.oldFilterColumnMap
	return nil
}
