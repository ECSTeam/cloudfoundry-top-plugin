package eventdata

import (
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type ContainerStats struct {
	ContainerMetric *events.ContainerMetric
	LastUpdate      time.Time
	OutCount        int64
	ErrCount        int64
}

func NewContainerStats() *ContainerStats {
	stats := &ContainerStats{}
	return stats
}
