package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type logEvent struct{ msg string }

type Standard struct {
	*logrus.Logger
}

func New() *Standard {
	l := &Standard{logrus.New()}
	fm := logrus.FieldMap{
		logrus.FieldKeyTime:  "time",
		logrus.FieldKeyLevel: "level_name",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyFunc:  "caller",
	}
	l.Formatter = &logrus.JSONFormatter{FieldMap: fm}
	l.Out = os.Stdout
	l.Level = logrus.InfoLevel
	return l
}

const serviceName = "power-plant-monitoring-system service"

var (
	httpServerMsg              = logEvent{"%s - Error initializing http status server: %s"}
	httpServerInitializedMsg   = logEvent{"%s - Initializing http status server: %s"}
	fatalMsg                   = logEvent{"%s - Fatal error: %s"}
	httpServerErrorMsg         = logEvent{"%s - HTTP server error %s"}
	receivedInterruptSignalMsg = logEvent{"%s - Received interrupt signal, stopping gracefully"}
	httpServerShutdownErrorMsg = logEvent{"%s - Error closing listeners, or context timeout: %s"}
	sqlConnectionErrorMsg      = logEvent{"%s - Error opening sql server: %s"}
	closerErrorMsg             = logEvent{"%s - Error closing resource: %s"}
	brokerWaiterMsg            = logEvent{"%s - Waiting for messages for %s ..."}
	brokerInitializingMsg      = logEvent{"%s - Initializing broker message"}
	fatalBrokerMsg             = logEvent{"%s - Fatal Broker error: %s"}
	brokerErrMsg               = logEvent{"%s - Broker error: %s"}
)

func (l *Standard) HTTPServerInitializationError(err error) {
	l.Errorf(httpServerMsg.msg, serviceName, err)
}

func (l *Standard) HTTPServerInitialization(text string) {
	l.Infof(httpServerInitializedMsg.msg, serviceName, text)
}

func (l *Standard) FatalError(err error) {
	l.Fatalf(fatalMsg.msg, serviceName, err)
}

func (l *Standard) HTTPServerShutdownError(err error) {
	l.Infof(httpServerShutdownErrorMsg.msg, serviceName, err)
}

func (l *Standard) HTTPServerError(err error) {
	l.Errorf(httpServerErrorMsg.msg, serviceName, err)
}

func (l *Standard) ReceivedInterruptSignal() {
	l.Infof(receivedInterruptSignalMsg.msg, serviceName)
}

func (l *Standard) SQLConnectionError(err error) {
	l.Errorf(sqlConnectionErrorMsg.msg, serviceName, err)
}

func (l *Standard) ShowCloserError(err error) {
	l.Errorf(closerErrorMsg.msg, serviceName, err)
}

func (l *Standard) BrokerWaiter(waiterName string) {
	l.Infof(brokerWaiterMsg.msg, serviceName, waiterName)
}

func (l *Standard) BrokerInitialization() {
	l.Infof(brokerInitializingMsg.msg, serviceName)
}

func (l *Standard) ShowFatalBrokerError(err error) {
	l.Errorf(fatalBrokerMsg.msg, serviceName, err)
}

func (l *Standard) ShowBrokerError(err error) {
	l.Errorf(brokerErrMsg.msg, serviceName, err)
}
