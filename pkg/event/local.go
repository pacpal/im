package event

type LocalEventBus struct {
	handlers map[string][]Handler
}

func NewLocalEventBus() *LocalEventBus {
	return &LocalEventBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *LocalEventBus) Subscribe(eventType string, handler Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *LocalEventBus) Publish(event Event) {
	handlers, ok := b.handlers[event.GetEventType()]
	if !ok {
		return
	}
	for _, h := range handlers {
		go h(event)
	}
}