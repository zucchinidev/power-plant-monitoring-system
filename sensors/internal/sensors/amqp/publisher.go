package amqp

import (
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
)

type DataSensorPublisher struct {
	broker *broker.Broker
	queue  *amqp.Queue
}

func NewDataSensorPublisher(b *broker.Broker, queueName string) (*DataSensorPublisher, error) {
	q, err := b.CreateQueue(queueName)
	if err != nil {
		return nil, err
	}
	return &DataSensorPublisher{broker: b, queue: q}, nil
}

func (s *DataSensorPublisher) Publish(b []byte) error {
	return s.broker.Channel().Publish(
		"",
		s.queue.Name,
		false,
		false,
		amqp.Publishing{Body: b},
	)
}

func (p *DataSensorPublisher) Queue() *amqp.Queue {
	return p.queue
}
