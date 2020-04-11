package sensors

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}

type MessagePublisher interface {
	Publish([]byte) error
}
