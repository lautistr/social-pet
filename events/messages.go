package events

import "github.com/lautistr/social-pet/models"

type Message interface {
	Type() string
}

type CreatedPostMessage struct {
	Post *models.Post
}

func (m CreatedPostMessage) Type() string {
	return "created_post"
}
