package uiCommon

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
)

const (
	MAX_SORT_COLUMNS = 5
	AscendingText    = "( " + UpArrow + " ascending )"
	DescendingText   = "( " + DownArrow + " descending )"
)

type EditSortView struct {
	*EditColumnViewAbs

	sortPosition   int
	sortColumns    []*SortColumn
	oldSortColumns []*SortColumn
}

func NewEditSortView(masterUI masterUIInterface.MasterUIInterface, name string, listWidget *ListWidget) *EditSortView {
	w := &EditSortView{EditColumnViewAbs: NewEditColumnViewAbs(masterUI, name, listWidget)}
	w.width = 55
	w.height = 14
	w.title = "Edit Sort"

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

	w.sortColumns = make([]*SortColumn, MAX_SORT_COLUMNS)
	for i, sc := range listWidget.sortColumns {
		w.sortColumns[i] = sc
	}

	// Save old sort for cancel
	w.oldSortColumns = make([]*SortColumn, len(listWidget.sortColumns))
	for i, sc := range listWidget.sortColumns {
		scClone := &SortColumn{
			id:          sc.id,
			reverseSort: sc.reverseSort,
		}
		w.oldSortColumns[i] = scClone
	}
	return w
}

func (w *EditSortView) initialLayoutCallback(g *gocui.Gui, v *gocui.View) error {

	if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.keyArrowDownAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.keyArrowUpAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keySpaceAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyDelete, gocui.ModNone, w.keyDeleteAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyBackspace, gocui.ModNone, w.keyDeleteAction); err != nil {
		return err
	}
	if err := g.SetKeybinding(w.name, gocui.KeyBackspace2, gocui.ModNone, w.keyDeleteAction); err != nil {
		return err
	}
	return nil
}

func (w *EditSortView) refreshDisplayCallback(g *gocui.Gui, v *gocui.View) error {

	v.Clear()
	fmt.Fprintln(v, " ")
	fmt.Fprintln(v, "  RIGHT or LEFT arrow - highlight sort column")
	fmt.Fprintln(v, "  DOWN or UP arrow - select sort position")
	fmt.Fprintln(v, "  SPACE - select column and toggle sort direction")
	fmt.Fprintln(v, "  DELETE - remove sort from position")
	fmt.Fprintln(v, "  ENTER - apply sort, ESC to cancel")
	fmt.Fprintln(v, "")

	for i, sc := range w.sortColumns {
		fmt.Fprintf(v, "    ")
		if w.sortPosition == i {
			fmt.Fprintf(v, util.REVERSE_WHITE)
		}
		displayName := "--none--"
		if sc != nil {
			sortDirection := AscendingText
			if sc.reverseSort {
				sortDirection = DescendingText
			}
			columnLabel := w.listWidget.columnMap[sc.id].label
			displayName = fmt.Sprintf("%-13v %v", columnLabel, sortDirection)
		}
		fmt.Fprintf(v, " Sort #%v: %v \n", i+1, displayName)
		if w.sortPosition == i {
			fmt.Fprintf(v, util.CLEAR)
		}
	}
	return nil
}

func (w *EditSortView) keyArrowDownAction(g *gocui.Gui, v *gocui.View) error {
	if w.sortColumns[w.sortPosition] == nil {
		return nil
	}
	if w.sortPosition+1 < MAX_SORT_COLUMNS {
		w.sortPosition++
	}
	return w.RefreshDisplay(g)
}

func (w *EditSortView) keyArrowUpAction(g *gocui.Gui, v *gocui.View) error {
	if w.sortPosition > 0 {
		w.sortPosition--
	}
	return w.RefreshDisplay(g)
}

func (w *EditSortView) keyDeleteAction(g *gocui.Gui, v *gocui.View) error {
	w.sortColumns[w.sortPosition] = nil
	pos := 0
	for _, sc := range w.sortColumns {
		if sc != nil {
			w.sortColumns[pos] = sc
			pos++
		}
	}
	for i := pos; i < len(w.sortColumns); i++ {
		w.sortColumns[i] = nil
	}
	w.applySort(g)
	return nil
}

func (w *EditSortView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

	sc := w.sortColumns[w.sortPosition]
	columnId := w.listWidget.selectedColumnId
	if sc == nil {
		sc = &SortColumn{
			id:          columnId,
			reverseSort: w.listWidget.columnMap[columnId].defaultReverseSort,
		}
		w.sortColumns[w.sortPosition] = sc
	} else {
		if sc.id == columnId {
			sc.reverseSort = !sc.reverseSort
		} else {
			sc.id = columnId
			sc.reverseSort = w.listWidget.columnMap[columnId].defaultReverseSort
		}

	}
	w.applySort(g)
	return nil
}

func (w *EditSortView) applySort(g *gocui.Gui) {

	useSortColumns := make([]*SortColumn, 0)
	for _, sc := range w.sortColumns {
		if sc != nil {
			useSortColumns = append(useSortColumns, sc)
		}
	}
	w.listWidget.sortColumns = useSortColumns
	w.listWidget.FilterAndSortData()
	w.listWidget.displayRowIndexOffset = 0
	w.RefreshDisplay(g)
}

func (w *EditSortView) applyActionCallback(g *gocui.Gui, v *gocui.View) error {
	w.applySort(g)
	return nil
}

func (w *EditSortView) cancelActionCallback(g *gocui.Gui, v *gocui.View) error {
	w.listWidget.sortColumns = w.oldSortColumns
	w.listWidget.FilterAndSortData()
	return nil
}
