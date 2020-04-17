package amqp

import (
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
)

type CoordinatorStartedReceptor struct {
	broker        *broker.Broker
	queue         *amqp.Queue
	nameEmitterFn func() error
}

func NewCoordinatorStartedReceptor(b *broker.Broker, nameEmitterFn func() error) (*CoordinatorStartedReceptor, error) {
	q, err := b.CreateQueue("queueName")
	if err != nil {
		return nil, err
	}
	coo := &CoordinatorStartedReceptor{broker: b, queue: q, nameEmitterFn: nameEmitterFn}
	err = coo.bind()
	if err != nil {
		return nil, err
	}
	return coo, nil
}

func (r *CoordinatorStartedReceptor) bind() error {
	return r.broker.Channel().QueueBind(
		r.queue.Name,
		"", // every events
		broker.SensorDiscoveryExchange,
		false,
		nil)
}

func (r *CoordinatorStartedReceptor) ListenForDiscoverCoordinators() {
	go func() {
		msgs, _ := r.broker.Channel().Consume(
			r.queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil)
		for range msgs {
			_ = r.nameEmitterFn()
		}
	}()
}
