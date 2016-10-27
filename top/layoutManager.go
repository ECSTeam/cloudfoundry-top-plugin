package top

import (
	//"fmt"
  "github.com/jroimartin/gocui"
)

type LayoutManager struct {
	 managers []gocui.Manager
}

func NewLayoutManager() *LayoutManager {
	return &LayoutManager{ }
}

func (w *LayoutManager) Layout(g *gocui.Gui) error {
  for _, m := range w.managers {
    if err := m.Layout(g); err != nil {
      return err
    }
  }
	return nil
}

func (w *LayoutManager) Contains(managerToFind gocui.Manager) bool {
  for _, m := range w.managers {
     if m == managerToFind {
         return true
     }
   }
   return false
}

func (w *LayoutManager) Add(m gocui.Manager) {
	w.managers = append(w.managers, m)
}


func (w *LayoutManager) Remove(managerToRemove gocui.Manager) {

  filteredManagers := []gocui.Manager{}
  for _, m := range w.managers {
     if m != managerToRemove {
         filteredManagers = append(filteredManagers, m)
     }
   }
   //fmt.Printf("\n\n\n[filteredManagers size:%v]", len(w.managers))
   w.managers = filteredManagers
}
