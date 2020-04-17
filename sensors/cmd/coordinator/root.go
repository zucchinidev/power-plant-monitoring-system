package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/internal/coordinator/amqp"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/broker"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/conf"
	"github.com/zucchinidev/power-plant-monitoring-system/sensors/shared/adapters/logger"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:   "coordinator",
	Short: "coordinator",
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

		go amqp.NewQueueDiscoverer(brokerManager).ListenForNewSource()

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
