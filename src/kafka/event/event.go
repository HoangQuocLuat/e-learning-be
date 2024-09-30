package kafka_event

type Event interface {
	GetId() string
}
