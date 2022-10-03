package main

import (
	"card-api/cmd/entities"
	"card-api/cmd/events"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var publisher = events.NewPublisher()

func helloWorld(w http.ResponseWriter, r *http.Request) {

	accountId, _ := uuid.NewUUID()

	account := entities.NewOwner(accountId)

	event := account.PublishNewCardCreated()

	publisher.Subscribe(event, &entities.BuildCardHandler{}, &entities.NotifyBuildingCardHandler{})

	go publisher.Notify(event)

	fmt.Println(event.Name())

	resp, _ := json.Marshal(*account)

	w.Write(resp)
}

func main() {
	http.HandleFunc("/api", helloWorld)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
