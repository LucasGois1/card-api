package entities

import (
	"card-api/cmd/events"
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"time"
)

type Card struct {
	Owner   uuid.UUID
	Number  int64
	DueDate time.Time
	Cvc     uint8
	Virtual bool
	Active  bool
}

type CardCreatedEvent struct {
	Identifier uuid.UUID
	Owner      uuid.UUID
	Number     int64
	DueDate    time.Time
	CreatedAt  time.Time
	Cvc        uint8
}

func (event *CardCreatedEvent) Name() string {
	return "CardCreatedEvent"
}

func (event *CardCreatedEvent) Id() uuid.UUID {
	return event.Identifier
}

func NewCard(uuid uuid.UUID) *Card {
	return &Card{
		uuid,
		int64(math.Abs(float64(rand.Int()))),
		time.Now().Add(time.Hour * 24 * 30 * 12 * 6).Local(),
		uint8(rand.Int()),
		false,
		false,
	}
}

type Credit struct {
	Owner   uuid.UUID
	DueDate time.Time
	Limit   float64
	Used    float64
	Card    []Card
}

type Debit struct {
	Owner   uuid.UUID
	Balance float64
	Card    []Card
}

type Owner struct {
	AccountId    uuid.UUID
	Credit       *Credit
	Debit        *Debit
	CardShipped  bool
	CardReceived bool
	CreatedAt    time.Time
}

func (o *Owner) PublishNewCardCreated() events.Event {
	cards := o.Debit.Card

	card := cards[0]
	event := CardCreatedEvent{
		uuid.New(),
		o.AccountId,
		card.Number,
		card.DueDate,
		time.Now(),
		card.Cvc,
	}

	return &event
}

func NewOwner(uuid uuid.UUID) *Owner {
	newCard := NewCard(uuid)

	listCard := []Card{*newCard}

	debitAccount := Debit{
		uuid,
		0.0,
		listCard,
	}

	owner := Owner{
		uuid,
		nil,
		&debitAccount,
		false,
		false,
		time.Now(),
	}

	return &owner
}

type BuildCardHandler struct{}

func (*BuildCardHandler) HandleEvent(event *chan *events.Event) {
	colorYellow := "\033[33m"

	for {
		log.Println(colorYellow, "BuildCardHandler:> Waiting for event...")
		e := <-*event

		log.Println(colorYellow, "BuildCardHandler:> New event received. Sending card information to the manufacturer...")
		log.Println(e)

		time.Sleep(time.Second * 5)

		log.Println(colorYellow, "BuildCardHandler:> Event processed.")
	}
}

type NotifyBuildingCardHandler struct{}

func (*NotifyBuildingCardHandler) HandleEvent(event *chan *events.Event) {
	colorBlue := "\033[34m"

	for {
		log.Println(colorBlue, "NotifyBuildingCardHandler:> Waiting for event...")
		e := <-*event

		log.Println(colorBlue, "NotifyBuildingCardHandler:> New event received. Sending email to the card owner 'You should receive your new card soon'")
		log.Println(colorBlue, e)

		time.Sleep(time.Second * 5)

		log.Println(colorBlue, "NotifyBuildingCardHandler:> Event processed.")
	}
}
