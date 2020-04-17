package amqp

import "time"

type EventAggregator struct {
	listeners map[string][]func(EventData)
}

func NewEventAggregator() *EventAggregator {
	return &EventAggregator{listeners: make(map[string][]func(EventData))}
}

func (ea *EventAggregator) AddListener(eventName string, f func(EventData)) {
	ea.listeners[eventName] = append(ea.listeners[eventName], f)
}

func (ea *EventAggregator) PublishEvent(eventName string, eventData EventData) {
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
