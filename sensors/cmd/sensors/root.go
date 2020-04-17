package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/sensors/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/sensors/sendingDataSensor"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/shared/message"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/conf"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/logger"
)

var (
	name  string
	freq  uint
	max   float64
	min   float64
	step  float64
	r     = rand.New(rand.NewSource(time.Now().UnixNano()))
	value float64
	nom   float64
)

func init() {
	rootCmd.Flags().StringVar(&name, "name", "sensor", "name of the sensor")
	rootCmd.Flags().UintVar(&freq, "freq", 5, "update frequency in cycles/sec")
	rootCmd.Flags().Float64Var(&max, "max", 5., "maximum value for generated readings")
	rootCmd.Flags().Float64Var(&min, "min", 1., "minimum value for generated readings")
	rootCmd.Flags().Float64Var(&step, "step", 0.1, "maximum allowable change per measurement")
	value = r.Float64()*(max-min) + min
	nom = (max-min)/2 + min
	gob.Register(message.SensorMessage{})
}

func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = step
		minStep = -1 * step * (value - min) / (nom - min)
	} else {
		maxStep = step * (max - value) / (max - nom)
		minStep = -1 * step
	}
	value += r.Float64()*(maxStep-minStep) + minStep
}

var rootCmd = &cobra.Command{
	Use:   "sensors",
	Short: "sensors",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			l             = logger.New()
			c             = conf.C()
			brokerManager *broker.Broker
			err           error
		)

		bCnf := broker.Cfg{Exchange: c.BrokerExchange, Url: c.BrokerQUrl}
		if brokerManager, err = initBrokerManager(bCnf, l); err != nil {
			l.ShowFatalBrokerError(err)
			return
		}
		defer brokerManager.Close()

		// we'll convert cycles/sec en milliseconds/cycles
		// e.g: 5 cycles/sec == 200 milliseconds/cycle
		dur, err := time.ParseDuration(strconv.Itoa(1000/int(freq)) + "ms")
		if err != nil {
			l.Panic(err)
		}
		signalTick := time.Tick(dur)
		buf := new(bytes.Buffer)
		queueNameEmitter := amqp.NewSensorQueueNameEmitter(brokerManager)
		publisher, err := amqp.NewPublisher(brokerManager, name, queueNameEmitter)
		if err != nil {
			l.Panic(err)
		}
		dataSensorSender := sendingDataSensor.NewService(publisher)

		go func() {
			for range signalTick {
				calcValue()
				reading := message.SensorMessage{
					Name:      name,
					Value:     value,
					Timestamp: time.Now().UTC(),
				}

				buf.Reset()
				if err := gob.NewEncoder(buf).Encode(reading); err != nil {
					l.Panic(err)
				}

				if err := dataSensorSender.Invoke(buf.Bytes()); err != nil {
					l.Panic(err)
				}

				log.Printf("Reading sent. Value: %v\n", value)
			}
		}()
		shutdown := make(chan struct{}, 1)
		killFn := terminate([]io.Closer{brokerManager}, l)
		go interruptSignal(l, shutdown)
		for {
			select {
			case err := <-brokerManager.Err:
				l.ShowBrokerError(err)
			case <-shutdown:
				killFn()
				return
			}
		}
	},
}

// Execute assembles the app commands necessaries to up the applications
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func terminate(closers []io.Closer, l *logger.Standard) func() {
	return func() {
		for _, closer := range closers {
			err := closer.Close()
			if err != nil {
				l.ShowCloserError(err)
			}
		}
	}
}

func interruptSignal(l *logger.Standard, shutdown chan struct{}) {
	signals := make(chan os.Signal, 1)
	// sigterm signal sent from kubernetes, interrupt signal sent from terminal
	signal.Notify(signals, syscall.SIGTERM, os.Interrupt)
	<-signals
	l.ReceivedInterruptSignal()
	shutdown <- struct{}{}
}

func initBrokerManager(brokerConf broker.Cfg, l *logger.Standard) (*broker.Broker, error) {
	brokerManager := broker.New(brokerConf)
	l.BrokerInitialization()
	if err := brokerManager.Connect(); err != nil {
		return nil, err
	}
	return brokerManager, nil
}
