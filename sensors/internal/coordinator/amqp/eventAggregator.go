package amqp

import "time"

type callable func(interface{})

type EventRaiser interface {
	AddListener(eventName string, f callable)
}

type EventAggregator struct {
	listeners map[string][]callable
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{listeners: make(map[string][]callable)}
}

func (ea *EventAggregator) AddListener(eventName string, f callable) {
	ea.listeners[eventName] = append(ea.listeners[eventName], f)
}

func (ea *EventAggregator) PublishEvent(eventName string, eventData interface{}) {
	if listeners, ok := ea.listeners[eventName]; ok {
		for _, listener := range listeners {
			listener(eventData)
		}
	}
}

type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

// Register allows register an event listener
