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
	"fmt"
	"strconv"
	"strings"
)

const (
	// Unicode characters: http://graphemica.com/unicode/characters/page/34
	Ellipsis = string('\U00002026')
)

func CaseInsensitiveLess(s1, s2 string) bool {
	// TODO: Find a more efficent way to do this that does not involve obj creation
	return strings.ToUpper(s1) < strings.ToUpper(s2)
}

func Format(n int64) string {
	return FormatUint64(uint64(n))
}

func FormatUint64(n uint64) string {
	in := strconv.FormatUint(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func FormatDisplayData(value string, size int) string {
	if len(value) > size {
		value = value[0:size-1] + Ellipsis
	}
	format := fmt.Sprintf("%%-%v.%vv", size, size)
	return fmt.Sprintf(format, value)
}
