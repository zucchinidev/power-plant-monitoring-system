package sendingDataSensor

import (
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/sensors"
)

type Service interface {
	Invoke(b []byte) error
}

type service struct {
	publisher sensors.MessagePublisher
}

func NewService(publisher sensors.MessagePublisher) Service {
	return &service{publisher: publisher}
}

func (s *service) Invoke(b []byte) error {
	return s.publisher.Publish(b)
}
