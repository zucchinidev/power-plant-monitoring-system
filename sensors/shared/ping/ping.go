package ping

type Pinger interface {
	Ping() error
}
