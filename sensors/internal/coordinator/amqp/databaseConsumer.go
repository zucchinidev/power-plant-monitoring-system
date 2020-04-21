package amqp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/shared/message"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
	"time"
)

const maxRate = 5 * time.Second

type DatabaseConsumer struct {
	broker      *broker.Broker
	queue       *amqp.Queue
	eventRaiser EventRaiser
	sources     []string
}

func NewDatabaseConsumer(broker *broker.Broker, eventRaiser EventRaiser) (*DatabaseConsumer, error) {
	q, err := broker.CreateQueue("PersistReadingsQueue")
	if err != nil {
		return nil, err
	}
	return &DatabaseConsumer{broker: broker, queue: q, eventRaiser: eventRaiser}, nil
}

func (d *DatabaseConsumer) StartConsumer() {
	d.eventRaiser.AddListener(DataSourceDiscovered, func(sensorName interface{}) {
		d.subscribeToDataEvent(sensorName.(string))
	})
}

func (d *DatabaseConsumer) subscribeToDataEvent(sensorName string) {
	for _, registerSource := range d.sources {
		if registerSource == sensorName {
			return // sensor already registered
		}
	}

	d.eventRaiser.AddListener(fmt.Sprintf("MessageReceived_%s", sensorName), func() func(interface{}) {
		prevTime := time.Unix(0, 0)
		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			ed := eventData.(EventData)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()
				sm := message.SensorMessage{
					Name:      ed.Name,
					Value:     ed.Value,
					Timestamp: ed.Timestamp,
				}

				buf.Reset() // reset previous values

				_ = gob.NewEncoder(buf).Encode(sm)
				_ = d.broker.Channel().Publish(
					"",
					d.queue.Name,
					false,
					false,
					amqp.Publishing{Body: buf.Bytes()})
			}

		}
	}())
}
