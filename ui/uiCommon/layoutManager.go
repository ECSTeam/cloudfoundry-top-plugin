package uiCommon

import (
	//"fmt"

	"log"

	"github.com/jroimartin/gocui"
	"github.com/kkellner/cloudfoundry-top-plugin/ui/masterUIInterface"
)

type LayoutManager struct {
	managers  []masterUIInterface.Manager
	viewNames []string
}

func NewLayoutManager() *LayoutManager {
	return &LayoutManager{}
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

func (w *LayoutManager) Add(addMgr masterUIInterface.Manager) {

	// TODO: This is just a development check -- can move later

	for _, m := range w.managers {
		if m.Name() == addMgr.Name() {
			log.Panicf("Attempting to add a ui manager named %v and it already exists", m.Name())
		}
	}

	w.managers = append(w.managers, addMgr)
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
		if m.Name() != managerToRemove.Name() {
			filteredManagers = append(filteredManagers, m)
		}
	}
	//fmt.Printf("\n\n\n[filteredManagers size:%v]", len(w.managers))
	w.managers = filteredManagers
	return w.Top()
}

func (w *LayoutManager) RemoveByName(managerViewNameToRemove string) masterUIInterface.Manager {
	m := w.GetManagerByViewName(managerViewNameToRemove)
	return w.Remove(m)
}

func (w *LayoutManager) GetManagerByViewName(managerViewNameToRemove string) masterUIInterface.Manager {
	for _, m := range w.managers {
		if m.Name() == managerViewNameToRemove {
			return m
		}
	}
	return nil
}
