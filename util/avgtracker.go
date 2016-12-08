package util

import (
	//"strconv"
	"time"
	//"fmt"
	"sync"
	//"gopkg.in/eapache/queue.v1"
	"github.com/eapache/queue"
	//"github.com/kkellner/cloudfoundry-top-plugin/toplog"
	//"fmt"
)

/*
func main() {
    //rateCounter := NewRateCounter(time.Minute)
    avgTracker := NewAvgTracker(time.Second)
    for i := 0; i < 9; i++ {
      avgTracker.Incr(100)
      time.Sleep(time.Millisecond * 50)
    }
		avgTracker.Incr(200)
    fmt.Printf("Rate: %v\n", avgTracker.Rate())
		fmt.Printf("Avg: %f\n", avgTracker.Avg())
}
*/

// A RateCounter is a thread-safe counter which returns the number of times
// 'Incr' has been called in the last interval
type AvgTracker struct {
	//counter  Counter
	interval   time.Duration
	timeQueue  *queue.Queue
	valueQueue *queue.Queue
	mu         *sync.Mutex
	totalValue int64
}

func AvgMultipleTrackers(trackers []*AvgTracker) float64 {

	totalLen := 0
	totalValue := int64(0)

	//fmt.Printf("\n\n***** len: [%v] **** \n", len(trackers))

	for _, tracker := range trackers {

		//fmt.Printf("\n\n***** interval: [%v] **** \n", tracker)
		//toplog.Debug(msg)
		//time.Sleep(10 * time.Second)

		tracker.mu.Lock()
		tracker.removeOld()
		totalLen = totalLen + tracker.valueQueue.Length()
		totalValue = totalValue + tracker.totalValue
		tracker.mu.Unlock()
	}

	avg := float64(-1)
	if totalLen > 0 {
		avg = float64(totalValue) / float64(totalLen)
	}
	return avg
}

// Constructs a new RateCounter, for the interval provided
func NewAvgTracker(intrvl time.Duration) *AvgTracker {
	return &AvgTracker{
		interval:   intrvl,
		timeQueue:  queue.New(),
		valueQueue: queue.New(),
		mu:         &sync.Mutex{},
	}
}

// Add an event into the RateCounter
func (r *AvgTracker) Track(val int64) {
	r.mu.Lock()
	r.removeOld()
	r.totalValue = r.totalValue + val
	r.timeQueue.Add(time.Now())
	r.valueQueue.Add(val)
	r.mu.Unlock()
}

// Return the current number of events in the last interval
func (r *AvgTracker) Rate() int {
	r.mu.Lock()
	r.removeOld()
	len := r.timeQueue.Length()
	r.mu.Unlock()
	return len
}

func (r *AvgTracker) Avg() float64 {
	r.mu.Lock()
	r.removeOld()
	len := r.valueQueue.Length()
	avg := float64(-1)
	if len > 0 {
		avg = float64(r.totalValue) / float64(len)
	}
	r.mu.Unlock()
	return avg
}

func (r *AvgTracker) removeOld() {

	if r.timeQueue.Length() > 0 {
		now := time.Now()
		for r.timeQueue.Length() > 0 {
			ts := r.timeQueue.Peek().(time.Time)
			if now.Sub(ts) < r.interval {
				break
			}
			//fmt.Printf("Remove - Now:[%v] ts:[%v] len:%v\n", now, ts, r.queue.Length())
			r.timeQueue.Remove()
			oldValue := r.valueQueue.Remove().(int64)
			r.totalValue = r.totalValue - oldValue
		}
	}

}
