package amqp

import (
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
)

const sensorListQueue = "SensorList"

type SensorListQueueTracker struct {
	broker *broker.Broker
	queue  *amqp.Queue
}

func NewSensorListQueueTracker(b *broker.Broker) (*SensorListQueueTracker, error) {
	q, err := b.CreateQueue(sensorListQueue)
	if err != nil {
		return nil, err
	}
	return &SensorListQueueTracker{broker: b, queue: q}, nil
}

func (s *SensorListQueueTracker) TrackQueue(q *amqp.Queue) error {
	return s.broker.Channel().Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{Body: []byte(q.Name)})
}
