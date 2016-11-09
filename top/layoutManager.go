package top

import (
	//"fmt"
  "github.com/jroimartin/gocui"
  "github.com/kkellner/cloudfoundry-top-plugin/masterUIInterface"
)

type LayoutManager struct {
	 managers []masterUIInterface.Manager
   viewNames []string
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

func (w *LayoutManager) Contains(managerToFind masterUIInterface.Manager) bool {
  for _, m := range w.managers {
     if m == managerToFind {
         return true
     }
   }
   return false
}

func (w *LayoutManager) Add(m masterUIInterface.Manager) {
	w.managers = append(w.managers, m)
}

func (w *LayoutManager) Top() masterUIInterface.Manager {
  len := len(w.managers)
  if len > 0 {
    return w.managers[len-1]
  }
  return nil
}



func (w *LayoutManager) Remove(managerToRemove masterUIInterface.Manager) masterUIInterface.Manager {

  filteredManagers := []masterUIInterface.Manager{}
  for _, m := range w.managers {
     if m != managerToRemove {
         filteredManagers = append(filteredManagers, m)
     }
   }
   //fmt.Printf("\n\n\n[filteredManagers size:%v]", len(w.managers))
   w.managers = filteredManagers
   return w.Top()
}
