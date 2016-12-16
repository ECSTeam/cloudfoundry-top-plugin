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

import "fmt"

const MEGABYTE = (1024 * 1024)
const GIGABYTE = (MEGABYTE * 1024)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	return b.StringWithPrecision(2)
}

func (b ByteSize) StringWithPrecision(precision int) string {

	format := fmt.Sprintf("%%.%vf", precision)

	switch {
	case b >= YB:
		return fmt.Sprintf(format+"YB", b/YB)
	case b >= ZB:
		return fmt.Sprintf(format+"ZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf(format+"EB", b/EB)
	case b >= PB:
		return fmt.Sprintf(format+"PB", b/PB)
	case b >= TB:
		return fmt.Sprintf(format+"TB", b/TB)
	case b >= GB:
		return fmt.Sprintf(format+"GB", b/GB)
	case b >= MB:
		return fmt.Sprintf(format+"MB", b/MB)
	case b >= KB:
		return fmt.Sprintf(format+"KB", b/KB)
	}
	return fmt.Sprintf(format+"B", b)
}

func main() {
	fmt.Println(YB, ByteSize(1e13))
}
