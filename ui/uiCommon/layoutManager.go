// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uiCommon

import (
	//"fmt"

	"log"

	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
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
		if m.Name() == managerToFind.Name() {
			return true
		}
	}
	return false
}

func (w *LayoutManager) ContainsViewName(viewName string) bool {
	for _, m := range w.managers {
		if m.Name() == viewName {
			return true
		}
	}
	return false
}

func (w *LayoutManager) SetCurrentView(viewName string) bool {
	// We remove then add the view back to get it to the bottom
	// of the list so that its considered "top" (or current)
	mgr := w.GetManagerByViewName(viewName)
	w.Remove(mgr)
	w.Add(mgr)
	return true
}

func (w *LayoutManager) AddToBack(addMgr masterUIInterface.Manager) {
	if w.Contains(addMgr) {
		log.Panicf("Attempting to add a ui manager named %v and it already exists", addMgr.Name())
	}
	w.managers = append([]masterUIInterface.Manager{addMgr}, w.managers...)
}

func (w *LayoutManager) Add(addMgr masterUIInterface.Manager) {
	if w.Contains(addMgr) {
		log.Panicf("Attempting to add a ui manager named %v and it already exists", addMgr.Name())
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
