package message

import "time"

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}
