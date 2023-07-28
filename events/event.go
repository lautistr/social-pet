package events

import (
	"github.com/lautistr/social-pet/models"
)

type EventStore interface {
	Close()
	PublishCreatedPost(post *models.Post) error
	SubscribeCreatedPost() (<-chan CreatedPostMessage, error)
	OnCreatedPost(f func(CreatedPostMessage)) error
}

var eventStore EventStore

func SetEventStore(store EventStore) {
	eventStore = store
}

func Close() {
	eventStore.Close()
}

func PublishCreatedPost(post *models.Post) error {
	return eventStore.PublishCreatedPost(post)
}

func SubscribeCreatedPost() (<-chan CreatedPostMessage, error) {
	return eventStore.SubscribeCreatedPost()
}

func OnCreatedPost(f func(CreatedPostMessage)) error {
	return eventStore.OnCreatedPost(f)
}
