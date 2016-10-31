package main

import (
	//"strconv"
	"time"
  //"gopkg.in/eapache/queue.v1"
  "github.com/eapache/queue"
)

func main() {
    queue := &NewRateCounter(time.Minute)
    queue.Incr()

}

// A RateCounter is a thread-safe counter which returns the number of times
// 'Incr' has been called in the last interval
type RateCounter struct {
	//counter  Counter
	interval time.Duration
  queue *queue.Queue
}

// Constructs a new RateCounter, for the interval provided
func NewRateCounter(intrvl time.Duration) *RateCounter {
	return &RateCounter{
		interval: intrvl,
    queue: &queue.Queue{},
	}
}

// Add an event into the RateCounter
func (r *RateCounter) Incr() {
	r.removeOld()
  r.queue.Add(time.Now())
}

// Return the current number of events in the last interval
func (r *RateCounter) Rate() int {
  r.removeOld()
	return r.queue.Length()
}

func (r *RateCounter) removeOld() {
  now := time.Now()
  for ts := r.queue.Peek().(time.Time); now.Sub(ts) > r.interval; {
    r.queue.Remove()
  }
}
