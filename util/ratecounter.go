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
	"sync"
	"time"

	"github.com/eapache/queue"
)

/*
func main() {
    //rateCounter := NewRateCounter(time.Minute)
    rateCounter := NewRateCounter(time.Second)
    for i := 0; i < 200; i++ {
      rateCounter.Incr()
      time.Sleep(time.Millisecond * 50)
    }
    fmt.Printf("%v\n", rateCounter.Rate())
}
*/

// A RateCounter is a thread-safe counter which returns the number of times
// 'Incr' has been called in the last interval
type RateCounter struct {
	interval time.Duration
	queue    *queue.Queue
	mu       sync.Mutex
}

// Constructs a new RateCounter, for the interval provided
func NewRateCounter(interval time.Duration) *RateCounter {
	return &RateCounter{
		interval: interval,
		queue:    queue.New(),
	}
}

// Add an event into the RateCounter
func (r *RateCounter) Incr() {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.removeOldFromTime(now)
	//fmt.Printf("now: %v\n", now)
	r.queue.Add(now)
}

// Return the current number of events in the last interval
func (r *RateCounter) Rate() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removeOld()
	len := r.queue.Length()
	return len
}

func (r *RateCounter) removeOld() {
	r.removeOldFromTime(time.Now())
}

func (r *RateCounter) removeOldFromTime(timeMark time.Time) {
	if r.queue.Length() > 0 {
		for r.queue.Length() > 0 {
			ts := r.queue.Peek().(time.Time)
			if timeMark.Sub(ts) < r.interval {
				break
			}
			r.queue.Remove()
		}
	}
}
