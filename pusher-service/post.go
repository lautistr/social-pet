package main

import (
	"time"

	"github.com/lautistr/social-pet/enums"
	"github.com/lautistr/social-pet/models"
)

type CreatedPostMessage struct {
	Type string `json:"type"`
	models.Post
}

func newCreatedPostMessage(id, body string, createdAt time.Time) *CreatedPostMessage {
	return &CreatedPostMessage{
		Type: enums.EventMessageType_Post,
		Post: models.Post{
			ID:        id,
			Body:      body,
			CreatedAt: createdAt,
		},
	}
}
