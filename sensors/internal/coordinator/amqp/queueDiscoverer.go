package amqp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/shared/message"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
	"log"
)

type QueueDiscoverer struct {
	broker *broker.Broker

	// With this map we might be able to close down listeners if the associated sensor goes off-line
	// (I'm not planning on implementing that),
	// we will able to pull the sensor and then publish their queue names
	// moreover we will able to avoid register the same sensor
	sources         map[string]<-chan amqp.Delivery
	eventAggregator *EventAggregator
}

func NewQueueDiscoverer(broker *broker.Broker) *QueueDiscoverer {
	return &QueueDiscoverer{
		broker:          broker,
		sources:         make(map[string]<-chan amqp.Delivery),
		eventAggregator: NewEventAggregator()}
}
func (q *QueueDiscoverer) ListenForNewSource() {
	queue, _ := q.broker.CreateQueue("") // RabbitMQ will create a unique name for it

	// By default when we create a queue, it is bound to the default exchange
	// For this reason, we need to rebind this queue to the fan-out exchange, in order to receive
	// the names of the sensor.
	_ = q.broker.Channel().QueueBind(
		queue.Name,
		"", // fan-out exchange just ignore this field
		"amq.fanout",
		false, // we don't want channel to be closed if the binding doesn't succeed. We know the exchange and queue exists at this point
		nil,
	)

	// in this point we are able to create a consumer which will receive the names of the queues.
	msgs, _ := q.broker.Channel().Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	q.discoverSensors()

	// when a message come in, it's going to indicate that a new sensor has come online and
	// is ready to send readings into the system. In order to receive those messages, we are going to use the
	// channel consume method to get access sensor's queue
	for sensorNameChan := range msgs {
		sensorName := string(sensorNameChan.Body)
		q.eventAggregator.PublishEvent("DataSourceDiscovered", sensorName)
		sensorDataChan, _ := q.broker.Channel().Consume(
			sensorName,
			"",
			true,
			false,
			false,
			false,
			nil)
		if _, ok := q.sources[sensorName]; !ok {
			q.sources[sensorName] = sensorDataChan

			go q.addListener(sensorDataChan)
		}
	}

}

func (q *QueueDiscoverer) addListener(dataChan <-chan amqp.Delivery) {
	for msg := range dataChan {
		decoder := gob.NewDecoder(bytes.NewReader(msg.Body))
		sensorMessage := new(message.SensorMessage)
		_ = decoder.Decode(sensorMessage)

		fmt.Printf("Received message: %v\n", sensorMessage)

		ed := EventData{
			Name:      sensorMessage.Name,
			Value:     sensorMessage.Value,
			Timestamp: sensorMessage.Timestamp,
		}
		sensorName := msg.RoutingKey
		q.eventAggregator.PublishEvent(
			fmt.Sprintf("MessageReceived_%s", sensorName),
			ed,
		)
	}
}

func (q *QueueDiscoverer) discoverSensors() {
	err := q.broker.Channel().ExchangeDeclare(
		broker.SensorDiscoveryExchange,
		"fanout",
		false,
		false,
		false, // You'll set this to true if you want this exchange to reject external publish request
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}
	_ = q.broker.Channel().Publish(
		broker.SensorDiscoveryExchange,
		"", // signal the sensors that we're looking for them
		false,
		false,
		amqp.Publishing{})
}
