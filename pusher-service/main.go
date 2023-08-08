package main

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/events"
	"github.com/lautistr/social-pet/models"
)

func main() {
	var cfg models.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	hub := NewHub()

	events.SetupNatsMessageBroker(cfg.NatsAddress, func(m events.CreatedPostMessage) {
		hub.Broadcast(newCreatedPostMessage(m.Post.ID, m.Post.Body, m.Post.CreatedAt), nil)
	})

	go hub.Run()
	http.HandleFunc(enums.ENDPOINTS_WEBSOCKET, hub.HandleWebSocket)
	err = http.ListenAndServe(enums.POSTS_PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
