package events

import (
	"context"

	"github.com/lautistr/social-pet/models"
)

type EventStore interface {
	Close()
	PublishCreatedPost(ctx context.Context, post *models.Post) error
	SubscribeCreatedPost(ctx context.Context) (<-chan CreatedPostMessage, error)
	OnCreatedPost(ctx context.Context, f func(CreatedPostMessage)) error
}

var eventStore EventStore

func SetEventStore(store EventStore) {
	eventStore = store
}

func Close() {
	eventStore.Close()
}

func PublishCreatedPost(ctx context.Context, post *models.Post) error {
	return eventStore.PublishCreatedPost(ctx, post)
}

func SubscribeCreatedPost(ctx context.Context) (<-chan CreatedPostMessage, error) {
	return eventStore.SubscribeCreatedPost(ctx)
}

func OnCreatedPost(ctx context.Context, f func(CreatedPostMessage)) error {
	return eventStore.OnCreatedPost(ctx, f)
}
