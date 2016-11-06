package masterUIInterface

import (
	"fmt"
  "log"
  "github.com/jroimartin/gocui"
)

type EditSortView struct {
  masterUI MasterUIInterface
	name string
  width int
  height int
  listWidget *ListWidget
}

func NewEditSortView(masterUI MasterUIInterface, name string, listWidget *ListWidget) *EditSortView {
	w := &EditSortView{masterUI: masterUI, name: name, listWidget: listWidget}
  w.width = 50
  w.height = 4
  return w
}

func (w *EditSortView) Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
	v, err := g.SetView(w.name, maxX/2-(w.width/2), maxY/2-(w.height/2), maxX/2+(w.width/2), maxY/2+(w.height/2))
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
    v.Title = "Edit Sort"
    v.Frame = true
    fmt.Fprintln(v, "Right Arrow or Left Arrow to select sort column, \npress SPACE to select column for sorting. \nPress ENTER to apply sort")
    if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeView); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.keyArrowRightAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.keyArrowLeftAction); err != nil {
      return err
    }
    if err := g.SetKeybinding(w.name, gocui.KeySpace, gocui.ModNone, w.keySpaceAction); err != nil {
      return err
    }
    if err := w.masterUI.SetCurrentViewOnTop(g, w.name); err != nil {
      log.Panicln(err)
    }

	}
	return nil
}

func (w *EditSortView) keyArrowRightAction(g *gocui.Gui, v *gocui.View) error {
  columnId := w.listWidget.editSortColumnId
  columns := w.listWidget.columns
  columnsLen := len(columns)
  for i, col := range columns {
    if col.id == columnId && i+1 < columnsLen {
      columnId = columns[i+1].id
      break
    }
  }
  writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
  w.listWidget.editSortColumnId = columnId
  w.listWidget.RefreshDisplay(g)
  return nil
}

func (w *EditSortView) keyArrowLeftAction(g *gocui.Gui, v *gocui.View) error {
  columnId := w.listWidget.editSortColumnId
  columns := w.listWidget.columns
  //columnsLen := len(columns)
  for i, col := range columns {
    if col.id == columnId && i > 0 {
      columnId = columns[i-1].id
      break
    }
  }
  writeFooter(g, fmt.Sprintf("\r columnId: %v", columnId) )
  w.listWidget.editSortColumnId = columnId
  w.listWidget.RefreshDisplay(g)
  return nil
}

func (w *EditSortView) keySpaceAction(g *gocui.Gui, v *gocui.View) error {

  sortColumns := make([]string,0)
  sortColumns = append(sortColumns, w.listWidget.editSortColumnId)
  w.listWidget.sortColumns = sortColumns
  w.listWidget.RefreshDisplay(g)
  return nil
}

func (w *EditSortView) closeView(g *gocui.Gui, v *gocui.View) error {
  w.listWidget.enableSortEdit(false)
  if err := w.masterUI.CloseView(w, w.name); err != nil {
    return err
  }
  w.listWidget.RefreshDisplay(g)
	return nil
}
