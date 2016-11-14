package masterUIInterface

import (
	"fmt"
	"regexp"

	"github.com/jroimartin/gocui"
)

type EditFilterView struct {
	*EditColumnViewAbs

	labelWidget Manager
	inputWidget Manager

	oldFilterColumnMap map[string]*FilterColumn

	editField bool
}

func NewEditFilterView(masterUI MasterUIInterface, name string, listWidget *ListWidget) *EditFilterView {
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

	return nil
}

func (w *EditFilterView) refreshDisplayCallback(g *gocui.Gui, v *gocui.View) error {

	v.Clear()
	fmt.Fprintln(v, " ")
	selectedColId := w.listWidget.selectedColumnId
	col := w.listWidget.columnMap[selectedColId]
	filter := w.listWidget.filterColumnMap[selectedColId]
	filterText := "--none--"
	if filter != nil {
		filterText = filter.filterText
	}

	fmt.Fprintf(v, " Column name: %v\n", col.label)
	fmt.Fprintf(v, " Filter: %v\n\n", filterText)

	if w.editField {
		fmt.Fprintf(v, "\n\n")
	} else {
		fmt.Fprintln(v, " RIGHT or LEFT arrow - select column")
		fmt.Fprintln(v, " SPACE - select column to edit")
		fmt.Fprintln(v, " ENTER - apply filter")
		fmt.Fprintln(v, "")
	}

	return nil
}

func (w *EditFilterView) applyValueCallback(g *gocui.Gui, v *gocui.View, mgr Manager, inputValue string) error {

	parentView, err := g.View(w.name)
	if err != nil {
		return err
	}

	compiledRegex, err := regexp.Compile(inputValue)
	if err != nil {
		fmt.Fprintf(parentView, "\r Error: %v", err)
		return nil
	}

	selectedColId := w.listWidget.selectedColumnId
	filter := &FilterColumn{filterText: inputValue, compiledRegex: compiledRegex}
	w.listWidget.filterColumnMap[selectedColId] = filter

	g.Cursor = false
	if err := w.masterUI.CloseView(w.labelWidget); err != nil {
		return err
	}
	if err := w.masterUI.CloseView(w.inputWidget); err != nil {
		return err
	}
	w.editField = false
	w.RefreshDisplay(g)
	return nil
}

func (w *EditFilterView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

	labelText := "Filter"
	maxLength := 10
	valueText := ""
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
