package events

import "github.com/google/uuid"

type Event interface {
	Name() string
	Id() uuid.UUID
}
