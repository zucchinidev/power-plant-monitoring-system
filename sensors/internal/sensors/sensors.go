package sensors

type MessagePublisher interface {
	Publish([]byte) error
}
