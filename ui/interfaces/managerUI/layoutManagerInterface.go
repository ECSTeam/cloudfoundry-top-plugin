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

package managerUI

import "github.com/jroimartin/gocui"

type LayoutManagerInterface interface {
	Contains(Manager) bool
	ContainsViewName(viewName string) bool
	Add(Manager)
	AddToBack(Manager)
	Remove(Manager) Manager
	Top() Manager
	GetManagerByViewName(viewName string) Manager
	RemoveByName(managerViewNameToRemove string) Manager
}

type Manager interface {
	Layout(*gocui.Gui) error
	Name() string
}
