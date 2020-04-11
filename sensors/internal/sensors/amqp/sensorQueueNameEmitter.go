package amqp

import (
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
)

type SensorQueueNameEmitter struct {
	broker *broker.Broker
}

func NewSensorQueueNameEmitter(b *broker.Broker) *SensorQueueNameEmitter {
	return &SensorQueueNameEmitter{broker: b}
}

func (s *SensorQueueNameEmitter) Emit(q *amqp.Queue) error {
	return s.broker.Channel().Publish(
		"amq.fanout",
		"",
		false,
		false,
		amqp.Publishing{Body: []byte(q.Name)})
}
