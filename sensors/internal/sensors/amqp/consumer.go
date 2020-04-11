package amqp

import (
	"errors"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/consumer"
)

type dataSensorConsumer struct {
	sensorName string
}

func (c *dataSensorConsumer) Suffix() string {
	return "on_" + c.sensorName + "_data_is_received"
}

func (c *dataSensorConsumer) Topics() []string {
	return []string{c.sensorName}
}

func (c *dataSensorConsumer) Run([]byte) error {
	return errors.New("eee")
}

func NewDataSensorConsumer(sensorName string) consumer.Consumer {
	return &dataSensorConsumer{sensorName: sensorName}
}
