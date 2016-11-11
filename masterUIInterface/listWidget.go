package masterUIInterface

import (
  "log"
  "fmt"
  "errors"
  "bytes"
  "strconv"
	"github.com/jroimartin/gocui"
  //"github.com/kkellner/cloudfoundry-top-plugin/debug"
  "github.com/kkellner/cloudfoundry-top-plugin/util"
)

type getRowKeyFunc func(index int) string
type getRowDisplayFunc func(index int, isSelected bool) string
type getDisplayHeaderFunc func() string

type getListSizeFunc func() int
type changeSelectionCallbackFunc func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool

type DisplayViewInterface interface {
  RefreshDisplay(g *gocui.Gui) error
  SetDisplayPaused(paused bool)
  GetDisplayPaused() bool
  SortData()
}

type ListColumn struct {
  id      string
  label    string
  size    int
  leftJustifyLabel bool
  sortFunc  util.LessFunc
  defaultReverseSort bool
  displayFunc getRowDisplayFunc
}

const LOCK_COLUMNS = 1

type ListWidget struct {
  masterUI MasterUIInterface
	name string
  topMargin int
  bottomMargin int

  Title string

  displayView DisplayViewInterface

  highlightKey string
  displayRowIndexOffset int
  displayColIndexOffset int

  GetRowKey  getRowKeyFunc
  PreRowDisplayFunc  getRowDisplayFunc
  GetListSize getListSizeFunc

  columns []*ListColumn
  columnMap map[string]*ListColumn

  editSort bool
  editSortColumnId string

  sortColumns []*SortColumn


}

type SortColumn struct {
  id  string
  reverseSort bool
}

func NewSortColumn(id string, reverseSort bool) *SortColumn {
  return &SortColumn{ id:id, reverseSort:reverseSort}
}

func NewListColumn(
  id, label string,
  size int,
  leftJustifyLabel bool,
  sortFunc  util.LessFunc,
  defaultReverseSort bool,
  displayFunc getRowDisplayFunc) *ListColumn {
  column := &ListColumn {
    id: id,
    label: label,
    size: size,
    leftJustifyLabel: leftJustifyLabel,
    sortFunc: sortFunc,
    defaultReverseSort: defaultReverseSort,
    displayFunc: displayFunc,
  }

  return column
}


func NewListWidget(masterUI MasterUIInterface, name string,
  topMargin, bottomMargin int, displayView DisplayViewInterface,
  columns []*ListColumn) *ListWidget {
  w := &ListWidget {
    masterUI: masterUI,
    name: name,
    topMargin: topMargin,
    bottomMargin: bottomMargin,
    displayView: displayView,
    columns: columns,
    columnMap: make(map[string]*ListColumn),
  }
  for _, col := range columns {
    w.columnMap[col.id] = col
  }

  return w
}

func (w *ListWidget) Name() string {
  return w.name
}

func (w *ListWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  bottom := maxY-w.bottomMargin
  if (w.topMargin >= bottom) {
    bottom = w.topMargin + 1
  }
	v, err := g.SetView(w.name, 0, w.topMargin, maxX-1, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
      return errors.New(w.name+" (ListWidget) layout error:" + err.Error())
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


    if err := g.SetKeybinding(w.name, 'o', gocui.ModNone, w.editSortAction); err != nil {
      log.Panicln(err)
    }

    if err := g.SetKeybinding(w.name, 'p', gocui.ModNone, w.toggleDisplayPauseAction); err != nil {
      log.Panicln(err)
    }

    if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone,
      func(g *gocui.Gui, v *gocui.View) error {
         w.highlightKey = ""
         w.RefreshDisplay(g)
         return nil
    }); err != nil {
      log.Panicln(err)
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (asUI *ListWidget) HighlightKey() string {
  return asUI.highlightKey
}

func (asUI *ListWidget) SetSortColumns(sortColumns []*SortColumn) {
  asUI.sortColumns = sortColumns
}

func (asUI *ListWidget) GetSortFunctions() []util.LessFunc {

  sortFunctions := make([]util.LessFunc, 0)
  for _, sortColumn := range asUI.sortColumns {
    sc := asUI.columnMap[sortColumn.id]
    sortFunc := sc.sortFunc
    if sortColumn.reverseSort {
      sortFunc = util.Reverse(sortFunc)
    }
    sortFunctions = append(sortFunctions, sortFunc)
  }
  return sortFunctions
}

func (asUI *ListWidget) SortData() {
  asUI.displayView.SortData()
}

func (asUI *ListWidget) RefreshDisplay(g *gocui.Gui) error {

  v, err := g.View(asUI.name)
  if err != nil {
		return err
	}
  _, maxY := v.Size()
  maxRows := maxY - 1

  v.Clear()
  listSize := asUI.GetListSize()

  if listSize>0 {
    stopRowIndex := maxRows + asUI.displayRowIndexOffset
    asUI.writeHeader(v)
    // Loop through all rows
    for i:=0;i<listSize && i<stopRowIndex;i++ {
      if i < asUI.displayRowIndexOffset {
        continue
      }
      asUI.writeRowData(v, i)
    }
  } else {
		fmt.Fprintf(v, " \n No data yet...")
  }

  return nil
}

func (asUI *ListWidget) writeRowData(v *gocui.View, rowIndex int) {
  isSelected := false
  if asUI.GetRowKey(rowIndex) == asUI.highlightKey {
    fmt.Fprintf(v, util.REVERSE_GREEN)
    isSelected = true
  }

  if asUI.PreRowDisplayFunc!=nil {
    fmt.Fprint(v, asUI.PreRowDisplayFunc(rowIndex, isSelected))
  }

  // Loop through all columns
  for colIndex, column := range asUI.columns {
    if colIndex >= LOCK_COLUMNS && colIndex < asUI.displayColIndexOffset+LOCK_COLUMNS {
      continue
    }
    fmt.Fprint(v, column.displayFunc(rowIndex, isSelected))
    fmt.Fprint(v," ")
  }
  fmt.Fprintf(v, "\n")
  fmt.Fprintf(v, util.CLEAR)
}

func (asUI *ListWidget) writeHeader(v *gocui.View) {

  // Loop through all columns (for headers)
  for colIndex, column := range asUI.columns {
    if colIndex >= LOCK_COLUMNS && colIndex < asUI.displayColIndexOffset+LOCK_COLUMNS {
      continue
    }
    editSortColumn := false
    if asUI.editSort && asUI.editSortColumnId == column.id {
      editSortColumn = true
      fmt.Fprintf(v, util.REVERSE_WHITE)
    }
    var buffer bytes.Buffer
    buffer.WriteString("%")
    if column.leftJustifyLabel {
      buffer.WriteString("-")
    }
    buffer.WriteString(strconv.Itoa(column.size))
    buffer.WriteString("v ")
    fmt.Fprintf(v, buffer.String(), column.label)
    if editSortColumn {
      fmt.Fprintf(v, util.CLEAR)
    }
  }
  fmt.Fprintf(v, "\n")
}

func (asUI *ListWidget) arrowRight(g *gocui.Gui, v *gocui.View) error {
  asUI.displayColIndexOffset++
  minColumnsToLeaveOnRight := 2
  maxOffset := len(asUI.columns) - (minColumnsToLeaveOnRight + LOCK_COLUMNS)
  if asUI.displayColIndexOffset > maxOffset {
    asUI.displayColIndexOffset = maxOffset
  }
  return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowLeft(g *gocui.Gui, v *gocui.View) error {
  asUI.displayColIndexOffset--
  if asUI.displayColIndexOffset < 0 {
    asUI.displayColIndexOffset = 0
  }
  return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) arrowUp(g *gocui.Gui, v *gocui.View) error {
  listSize := asUI.GetListSize()
  callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
    if rowIndex > 0 {
      _, viewY := v.Size()
      viewSize := viewY-1
      asUI.highlightKey = lastKey
      offset := rowIndex-1
      if (offset > listSize - viewSize) {
        offset = listSize - viewSize
      }
      if (asUI.displayRowIndexOffset > offset || rowIndex > asUI.displayRowIndexOffset+viewSize) {
        asUI.displayRowIndexOffset = offset
      }
      return true
    }
    return false
  }
  return asUI.moveHighlight(g, v, callbackFunc)
}

func (asUI *ListWidget) arrowDown(g *gocui.Gui, v *gocui.View) error {

  listSize := asUI.GetListSize()
  callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
    if rowIndex+1 < listSize {
      _, viewY := v.Size()
      offset := (rowIndex+2) - (viewY-1)
      if (offset>asUI.displayRowIndexOffset || rowIndex < asUI.displayRowIndexOffset) {
        asUI.displayRowIndexOffset = offset
      }
      asUI.highlightKey = asUI.GetRowKey(rowIndex+1)
      return true
    }
    return false
  }
  return asUI.moveHighlight(g, v, callbackFunc)
}


func (asUI *ListWidget) moveHighlight(g *gocui.Gui, v *gocui.View, callback changeSelectionCallbackFunc) error {

  listSize := asUI.GetListSize()
  if asUI.highlightKey == "" {
    if listSize > 0 {
      asUI.highlightKey = asUI.GetRowKey(0)
    }
  } else {
    lastKey := ""
    for rowIndex:=0;rowIndex<listSize;rowIndex++ {
      if asUI.GetRowKey(rowIndex) == asUI.highlightKey {
        if callback(g,v,rowIndex,lastKey) {
          break
        }
      }
      lastKey = asUI.GetRowKey(rowIndex)
    }
  }
  return asUI.RefreshDisplay(g)

}


func (asUI *ListWidget) pageUpAction(g *gocui.Gui, v *gocui.View) error {
  callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
    if rowIndex > 0 {
      _, viewY := v.Size()
      viewSize := viewY-1
      offset := 0
      if rowIndex == asUI.displayRowIndexOffset {
        offset = rowIndex - viewSize
      } else {
        offset = rowIndex - (rowIndex - asUI.displayRowIndexOffset)
      }
      if offset < 0 {
        offset = 0
      }
      if (offset < asUI.displayRowIndexOffset) {
        asUI.displayRowIndexOffset = offset
      }
      asUI.highlightKey = asUI.GetRowKey(offset)
      return true
    }
    return false
  }
  return asUI.moveHighlight(g, v, callbackFunc)
}

func (asUI *ListWidget) pageDownAction(g *gocui.Gui, v *gocui.View) error {
  listSize := asUI.GetListSize()
  callbackFunc := func(g *gocui.Gui, v *gocui.View, rowIndex int, lastKey string) bool {
    if rowIndex+1 < listSize {
      _, viewY := v.Size()
      viewSize := viewY-1
      offset := 0
      if rowIndex == (asUI.displayRowIndexOffset + viewSize - 1) {
        offset = rowIndex + viewSize
      } else {
        offset = rowIndex + (viewSize - (rowIndex - asUI.displayRowIndexOffset)) - 1
      }
      if offset > listSize-1 {
        offset = listSize - 1
      }
      if (offset > (asUI.displayRowIndexOffset + viewSize)) {
        asUI.displayRowIndexOffset = offset - viewSize + 1
      }
      asUI.highlightKey = asUI.GetRowKey(offset)
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

func (asUI *ListWidget) toggleDisplayPauseAction(g *gocui.Gui, v *gocui.View) error {

  asUI.displayView.SetDisplayPaused(!asUI.displayView.GetDisplayPaused())
  return asUI.displayView.RefreshDisplay(g)

}

func (asUI *ListWidget) editSortAction(g *gocui.Gui, v *gocui.View) error {

  asUI.editSort = true
  if asUI.editSortColumnId == "" {
    asUI.editSortColumnId = asUI.columns[0].id
  }

  editView := NewEditSortView(asUI.masterUI, asUI.name+".editView", asUI)
  asUI.masterUI.LayoutManager().Add(editView)
  asUI.masterUI.SetCurrentViewOnTop(g,"editView")
  return asUI.RefreshDisplay(g)
}

func (asUI *ListWidget) enableSortEdit(enable bool) {
  asUI.editSort = enable
}
