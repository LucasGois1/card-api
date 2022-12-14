package main

import (
	"card-api/cmd/entities"
	"card-api/cmd/events"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
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

const PORT = "8000"

func main() {
	args := os.Args[1:]

	port := ":"
	if len(args) != 1 {
		log.Println("Port not informed, using default value ", PORT)
		port += PORT
	} else {
		port += args[0]
	}

	http.HandleFunc("/api", helloWorld)
	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
