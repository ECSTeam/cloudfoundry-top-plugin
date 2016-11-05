package masterUIInterface

import (
  "log"
  "fmt"
  "errors"
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
  GetRowDisplay  getRowDisplayFunc
  GetDisplayHeader getDisplayHeaderFunc
  GetListSize getListSizeFunc


}

func NewListWidget(masterUI MasterUIInterface, name string,
  topMargin, bottomMargin int, displayView DisplayViewInterface,
  //displayItems []DisplayItemInterface
  ) *ListWidget {
  w := &ListWidget {
    masterUI: masterUI,
    name: name,
    topMargin: topMargin,
    bottomMargin: bottomMargin,
    displayView: displayView,
    //displayItems: displayItems,
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
  if asUI.GetRowDisplay==nil {
    debug.Debug(fmt.Sprintf("GetRowDisplay function not set"))
    return errors.New("GetRowDisplay function not set")
  }

  v, err := g.View(asUI.name)
  if err != nil {
		return err
	}
  _, maxY := v.Size()
  maxRows := maxY - 1

  v.Clear()
  listSize := asUI.GetListSize()

  //fmt.Fprintf(v, "list size: %v\n",listSize)
  fmt.Fprintln(v, asUI.GetDisplayHeader())
  //fmt.Fprintf(v, "\n")

  for i:=0;i<listSize && i<maxRows;i++ {
    isSelected := false
    if asUI.GetRowKey(i) == asUI.highlightKey {
      fmt.Fprintf(v, util.GREEN + util.REVERSE)
      isSelected = true
    }
    fmt.Fprintln(v, asUI.GetRowDisplay(i,isSelected))
    //fmt.Fprintf(v, "\n")
    fmt.Fprintf(v, util.CLEAR)
  }

  if listSize==0 {
		fmt.Fprintln(v, "No data yet...")
  }

  return nil
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

  asUI.displayView.RefreshDisplay(g)
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

  asUI.displayView.RefreshDisplay(g)
  return nil
}
