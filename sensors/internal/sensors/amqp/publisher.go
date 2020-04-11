package amqp

import (
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
)

type Publisher struct {
	broker *broker.Broker
	queue  *amqp.Queue
}

func NewPublisher(b *broker.Broker, queueName string, tracker *SensorListQueueTracker) (*Publisher, error) {
	q, err := b.CreateQueue(queueName)
	if err != nil {
		return nil, err
	}
	if err := tracker.TrackQueue(q); err != nil {
		return nil, err
	}
	return &Publisher{broker: b, queue: q}, nil
}

func (s *Publisher) Publish(b []byte) error {
	return s.broker.Channel().Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{Body: b},
	)
}
