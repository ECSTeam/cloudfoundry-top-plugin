package eventdata

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type ContainerStats struct {
	ContainerIndex  int
	Ip              string
	ContainerMetric *events.ContainerMetric
	LastUpdate      time.Time
	OutCount        int64
	ErrCount        int64
}

func NewContainerStats(containerIndex int) *ContainerStats {
	stats := &ContainerStats{ContainerIndex: containerIndex}
	return stats
}
