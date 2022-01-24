package pubsub

type Event struct {
	Id int
}

type PubSub interface {
	// publishes the event e to all current subscriptions.
	Publish (e Event)

	// registers c to receive future events.
	// if Publish(e1) happens before Publish(e2),
	// subscribers receive e1 before e2.
	Subscribe (c chan<- Event)

	// cancels the prior subscription of channel c.
	Cancel(c chan<- Event)
}
