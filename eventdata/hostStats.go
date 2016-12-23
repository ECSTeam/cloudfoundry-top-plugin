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

package eventdata

import "sort"
import "strings"

type HostSlice []*HostStats

type HostStats struct {
	// Key: path (needs to include empty string as key for root path)
	RouteStatsMap map[string]*RouteStats

	// If path is not found, dynamically register it to this menu levels deep in path
	// TODO: Problem -- if the first call is to "/v2" and its dynamically registered
	// then a subsequent call to /v2/apps" will to be added as there will be a match.
	// However, we really do want to register /v2/apps
	dynamicAddPathDepth int

	// index of paths where the best match is first (longest path first)
	pathIndex []string
}

func NewHostStats(hostName string) *HostStats {
	stats := &HostStats{}
	stats.RouteStatsMap = make(map[string]*RouteStats)
	return stats
}

func (hs *HostStats) AddPath(path string, routeId string) *RouteStats {
	rs := NewRouteStats(routeId)
	hs.RouteStatsMap[path] = rs
	hs.rebuildPathIndex()
	return rs
}

func (hs *HostStats) AddPathDynamic(fullPath string, routeId string) *RouteStats {
	// TODO: based on fullPath and hs.dynamicAddPathDepth tuncate the
	// fullPath if needed to get it to dynamicAddPathDepth size, then add
	rs := hs.AddPath(fullPath, routeId)
	return rs
}

// Build index of paths where the best match is first
func (hs *HostStats) rebuildPathIndex() {

	paths := make([]string, 0, len(hs.RouteStatsMap))
	for path, _ := range hs.RouteStatsMap {
		paths = append(paths, path)
	}
	sort.Sort(sort.Reverse(ByLength(paths)))
	hs.pathIndex = paths
}

// Find matching route for given path
// TODO: Are "path" definitions in CF case sensative?
//
// [0] "/webappa/subapp1"
// [1] "/webappa"
// [2] ""
//
// findPath = "/webappabc"    => ""
// findPath = "/webappa"	  => "/webappa"
// findPath = "/webappa/"	  => "/webappa"
// findPath = "/webappa/doc"  => "/webappa"
//
func (hs *HostStats) FindPathMatch(findPath string) string {

	// TODO: need to make sure we take into account dynamicAddPathDepth
	// e.g., do not return "/v2" if calling with "/v2/app" even if /v2
	// is registered and /v2/apps is not -- we should return:
	//		 empty match?  or "/v2/apps" even though its not in list?

	for _, path := range hs.pathIndex {
		if strings.HasPrefix(findPath, path) {
			pathLen := len(path)
			if len(findPath) == pathLen {
				return path
			}
			if findPath[pathLen] == '/' {
				return path
			}
		}
	}
	return ""
}

func (hs *HostStats) FindRouteStats(findPath string) *RouteStats {
	pathMatch := hs.FindPathMatch(findPath)
	return hs.RouteStatsMap[pathMatch]
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
