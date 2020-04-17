package broker

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/consumer"
)

const (
	prefix                  = "power_plant_monitoring_system"
	SensorDiscoveryExchange = "SensorDiscoveryExchange"
)

type Cfg struct {
	Exchange string
	Url      string
}

type Broker struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	c      Cfg
	Err    chan error
	queues []*amqp.Queue
}

func (b *Broker) Close() error {
	_ = b.ch.Close()
	return b.conn.Close()
}

func (b *Broker) Ping() error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	return ch.Close()
}

func (b *Broker) Listen(consumers []consumer.Consumer) error {
	for _, handler := range consumers {
		ch, err := b.conn.Channel()
		if err != nil {
			return err
		}

		queueName := fmt.Sprintf("%s_%s", prefix, handler.Suffix())
		q, err := b.bindTopics(ch, queueName, handler.Topics())
		if err != nil {
			return err
		}

		consumerName := fmt.Sprintf("%s_%s", prefix, handler.Suffix())
		delivery, err := ch.Consume(q.Name, consumerName, false, false, false, false, nil)
		if err != nil {
			return err
		}

		go func(delivery <-chan amqp.Delivery, consumer consumer.Consumer) {

			for d := range delivery {
				if err := consumer.Run(d.Body); err != nil {
					b.Err <- fmt.Errorf(consumer.Suffix()+" error (run) %s", err)
					if err := d.Nack(false, true); err != nil {
						b.Err <- fmt.Errorf(consumer.Suffix()+" error (nack) %s", err)
					}
					continue
				} else {
					if err := d.Ack(false); err != nil {
						b.Err <- fmt.Errorf(consumer.Suffix()+" error (ack) %s", err)
						continue
					}
				}
			}
		}(delivery, handler)
	}
	return nil
}

func (b *Broker) bindTopics(ch *amqp.Channel, queueName string, topics []string) (*amqp.Queue, error) {
	q, err := b.declareQueue(ch, queueName)
	if err != nil {
		return nil, err
	}

	for _, t := range topics {
		if err := ch.QueueBind(queueName, t, b.c.Exchange, false, amqp.Table{}); err != nil {
			return nil, err
		}
	}
	return &q, nil
}

func (b *Broker) declareQueue(channel *amqp.Channel, name string) (amqp.Queue, error) {
	return channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (b *Broker) Connect() error {
	var err error
	if b.conn, err = amqp.Dial(b.c.Url); err != nil {
		return err
	}

	if b.ch, err = b.conn.Channel(); err != nil {
		return err
	}
	return nil
}

func (b *Broker) Publish(routingKey string, body []byte) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(
		b.c.Exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Priority:     0,
		},
	)
}

func (b *Broker) CreateQueue(queueName string) (*amqp.Queue, error) {
	q, err := b.declareQueue(b.ch, queueName)
	if err != nil {
		return nil, err
	}
	b.queues = append(b.queues, &q)
	return &q, nil
}

func (b *Broker) Channel() *amqp.Channel {
	return b.ch
}

func New(c Cfg) *Broker {
	return &Broker{c: c, Err: make(chan error)}
}
