package consumer

type Consumer interface {
	Suffix() string
	Topics() []string
	Run([]byte) error
}
