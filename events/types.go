package events

const (
	Unknown Type = iota
	Message
	AskForLocation
	ReceiveLocation
	AskForPersistentLocation
)

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(event Event) error
}

type Type int

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
