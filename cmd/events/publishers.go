package events

type EventRegister struct {
	EventName           string
	Channel             *chan *Event
	listenersRegistered uint8
}

func NewEventRegister(name string) *EventRegister {
	return &EventRegister{
		EventName:           name,
		listenersRegistered: 0,
	}
}

func (e *EventRegister) IncreaseListeners() {
	e.listenersRegistered++
}

func (e *EventRegister) RegisterNewChannel() *chan *Event {
	channel := make(chan *Event)
	e.Channel = &channel

	return &channel
}

type Publisher struct {
	events    map[string]EventRegister
	listeners map[Listener]EventRegister
}

func NewPublisher() *Publisher {
	return &Publisher{
		events:    make(map[string]EventRegister),
		listeners: make(map[Listener]EventRegister),
	}
}

func (p *Publisher) subscribeEvent(event Event) *EventRegister {
	var registeredEvent *EventRegister

	if e, ok := p.events[event.Name()]; ok {
		registeredEvent = &e
	} else {
		registeredEvent = NewEventRegister(event.Name())
		registeredEvent.RegisterNewChannel()

		p.events[event.Name()] = *registeredEvent
	}

	return registeredEvent
}

func (p *Publisher) Subscribe(event Event, listener ...Listener) {
	eventInformation := p.subscribeEvent(event)

	for _, l := range listener {
		if e, ok := p.listeners[l]; !ok {

			if e.EventName == "" {
				eventInformation.IncreaseListeners()

				// update map information, can't update map values by reference
				p.listeners[l] = *eventInformation
				p.events[eventInformation.EventName] = *eventInformation

				go l.HandleEvent(eventInformation.Channel)
			}
		}
	}
}

func (p *Publisher) Unsubscribe(listener Listener) {
	// close all channels and remove listeners from map
	delete(p.listeners, listener)
}

func (p *Publisher) Notify(event ...Event) {
	// Send event to all respective channels
	for _, eventRegisters := range p.listeners {
		for _, e := range event {
			if eventRegisters.EventName == e.Name() {
				*eventRegisters.Channel <- &e
			}
		}
	}
}

type Listener interface {
	HandleEvent(event *chan *Event)
}
