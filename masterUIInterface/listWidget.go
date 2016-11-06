package masterUIInterface

import (
  "log"
  "fmt"
  "errors"
  "bytes"
  "strconv"
	"github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
  "github.com/kkellner/cloudfoundry-top-plugin/util"
)

type getRowKeyFunc func(index int) string
type getRowDisplayFunc func(index int, isSelected bool) string
type getDisplayHeaderFunc func() string

type getListSizeFunc func() int

type DisplayViewInterface interface {
  RefreshDisplay(g *gocui.Gui) error
}

type ListColumn struct {
  id      string
  label    string
  size    int
  leftJustifyLabel bool
  sortFunc  util.LessFunc
  reverseSort bool
  displayFunc getRowDisplayFunc
}

type ListWidget struct {
  masterUI MasterUIInterface
	name string
  topMargin int
  bottomMargin int

  Title string

  displayView DisplayViewInterface
  //displayItems []DisplayItemInterface

  highlightKey string
  displayIndexOffset int

  GetRowKey  getRowKeyFunc
  PreRowDisplayFunc  getRowDisplayFunc
  GetListSize getListSizeFunc

  columns []*ListColumn

  editSort bool
  editSortColumnId string

  sortColumns []string

}

func NewListColumn(
  id, label string,
  size int,
  leftJustifyLabel bool,
  sortFunc  util.LessFunc,
  reverseSort bool,
  displayFunc getRowDisplayFunc) *ListColumn {
  column := &ListColumn {
    id: id,
    label: label,
    size: size,
    leftJustifyLabel: leftJustifyLabel,
    sortFunc: sortFunc,
    reverseSort: reverseSort,
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
  }
  return w
}

func (w *ListWidget) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()

  //0, asUI.topMargin, maxX-1, maxY-asUI.bottomMargin

	v, err := g.SetView(w.name, 0, w.topMargin, maxX-1, maxY-w.bottomMargin)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = w.Title
    v.Frame = true
    //fmt.Fprintln(v, "Future home of filter screen")
    //if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeFilterWidget); err != nil {
    //  return err
    //}

    if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
      log.Panicln(err)
    }
    if err := g.SetKeybinding(w.name, 's', gocui.ModNone, w.editSortAction); err != nil {
      log.Panicln(err)
    }

    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (asUI *ListWidget) Name() string {
  return asUI.name
}

func (asUI *ListWidget) HighlightKey() string {
  return asUI.highlightKey
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
    asUI.writeHeader(v)
    for i:=0;i<listSize && i<maxRows;i++ {
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
    fmt.Fprintf(v, util.GREEN + util.REVERSE)
    isSelected = true
  }

  if asUI.PreRowDisplayFunc!=nil {
    fmt.Fprint(v, asUI.PreRowDisplayFunc(rowIndex, isSelected))
  }

  for _, column := range asUI.columns {
    fmt.Fprint(v, column.displayFunc(rowIndex, isSelected))
    fmt.Fprint(v," ")
  }
  fmt.Fprintf(v, "\n")
  fmt.Fprintf(v, util.CLEAR)
}


func (asUI *ListWidget) writeHeader(v *gocui.View) {


  for _, column := range asUI.columns {
    editSortColumn := false
    if asUI.editSort && asUI.editSortColumnId == column.id {
      editSortColumn = true
      fmt.Fprintf(v, util.WHITE + util.REVERSE)
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

func (asUI *ListWidget) arrowUp(g *gocui.Gui, v *gocui.View) error {

  if asUI.GetRowKey==nil {
    debug.Debug(fmt.Sprintf("GetRowKey function not set"))
    return errors.New("GetRowKey function not set")
  }

  listSize := asUI.GetListSize()
  if asUI.highlightKey == "" {
    if listSize > 0 {
      asUI.highlightKey = asUI.GetRowKey(0)
    }
  } else {
    lastKey := ""
    for i:=0;i<listSize;i++ {
      if asUI.GetRowKey(i) == asUI.highlightKey {
        if i > 0 {
          asUI.highlightKey = lastKey
          offset := i-1
          //writeFooter(g,"\r row["+strconv.Itoa(row)+"]")
          //writeFooter(g,"o["+strconv.Itoa(offset)+"]")
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayIndexOffset)+"]")
          if (asUI.displayIndexOffset > offset) {
            asUI.displayIndexOffset = offset
          }
          break
        }
      }
      lastKey = asUI.GetRowKey(i)
    }
  }

  asUI.RefreshDisplay(g)
  return nil

}

func (asUI *ListWidget) arrowDown(g *gocui.Gui, v *gocui.View) error {

  listSize := asUI.GetListSize()
  if asUI.highlightKey == "" {
    if listSize > 0 {
      asUI.highlightKey = asUI.GetRowKey(0)
    }
  } else {
    for i:=0;i<listSize;i++ {
      if asUI.GetRowKey(i) == asUI.highlightKey {
        if i+1 < listSize {
          asUI.highlightKey = asUI.GetRowKey(i+1)
          _, viewY := v.Size()
          offset := (i+2) - (viewY-1)
          if (offset>asUI.displayIndexOffset) {
            asUI.displayIndexOffset = offset
          }
          //writeFooter(g,"\r row["+strconv.Itoa(row)+"]")
          //writeFooter(g,"viewY["+strconv.Itoa(viewY)+"]")
          //writeFooter(g,"o["+strconv.Itoa(offset)+"]")
          //writeFooter(g,"rowOff["+strconv.Itoa(asUI.displayIndexOffset)+"]")
          break
        }
      }
    }
  }

  asUI.RefreshDisplay(g)
  return nil
}

func (asUI *ListWidget) RefreshDisplayX(g *gocui.Gui) {
  asUI.displayView.RefreshDisplay(g)
}


// This is for debugging -- remove it later
func writeFooter(g *gocui.Gui, msg string) {
  v, _ := g.View("footerView")
  fmt.Fprint(v, msg)

}

func (asUI *ListWidget) editSortAction(g *gocui.Gui, v *gocui.View) error {

  // TODO Freeze display update
  asUI.editSort = true
  if asUI.editSortColumnId == "" {
    asUI.editSortColumnId = asUI.columns[0].id
  }
  writeFooter(g,"\r EDIT SORT")

  editView := NewEditSortView(asUI.masterUI, "editView", asUI)
  asUI.masterUI.LayoutManager().Add(editView)
  asUI.masterUI.SetCurrentViewOnTop(g,"editView")


  asUI.RefreshDisplay(g)
  return nil

}

func (asUI *ListWidget) enableSortEdit(enable bool) {
  asUI.editSort = enable
}
