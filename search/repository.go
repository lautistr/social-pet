package search

import (
	"context"

	"github.com/lautistr/social-pet/models"
)

type SearchRepository interface {
	Close()
	IndexPost(ctx context.Context, post models.Post) error
	SearchPost(ctx context.Context, query string) ([]models.Post, error)
}

var repo SearchRepository

func SetSearchRepository(r SearchRepository) {
	repo = r
}

func Close() {
	repo.Close()
}

func IndexPost(ctx context.Context, post models.Post) error {
	return repo.IndexPost(ctx, post)
}

func SearchPost(ctx context.Context, query string) ([]models.Post, error) {
	return repo.SearchPost(ctx, query)
}
