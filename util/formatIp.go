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
	"hash/fnv"
	"net"
)

func Ip2long(s string) (ret int64) {
	bip := ([]byte)(net.ParseIP(s).To4())
	if len(bip) != 4 {
		// We don't have a IPv4 so give it a hash so at least we have a predictable sort order
		return int64(hash(s))
	}
	return (int64)(bip[0])*(1<<24) + (int64)(bip[1])*(1<<16) + (int64)(bip[2])*(1<<8) + (int64)(bip[3])
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
