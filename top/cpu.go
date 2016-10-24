package top

//#include <time.h>
import "C"
import "time"

var startTime = time.Now()
var startTicks = C.clock()

func CpuUsagePercent() float64 {
    clockSeconds := float64(C.clock()-startTicks) / float64(C.CLOCKS_PER_SEC)
    realSeconds := time.Since(startTime).Seconds()
    return clockSeconds / realSeconds * 100
}

func CpuResetAvg() {
  startTime = time.Now()
  startTicks = C.clock()
}
