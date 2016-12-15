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

package util

import (
    "sort"
)

type Sortable interface { }

type LessFunc func(p1, p2 Sortable) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	slice []Sortable
	less  []LessFunc
}

func Reverse(less LessFunc) LessFunc {
  return func(p1, p2 Sortable) bool { return less(p2, p1) }
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
//func (ms *multiSorter) Sort(appStats []*AppStats) {
func (ms *multiSorter) Sort(slice []Sortable) {
	ms.slice = slice
	//sort.Sort(sort.Reverse(ms))
  sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less []LessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.slice)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.slice[i], ms.slice[j] = ms.slice[j], ms.slice[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.slice[i], ms.slice[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
  if len(ms.less) > 0 {
	 return ms.less[k](p, q)
  }
  return false
}
