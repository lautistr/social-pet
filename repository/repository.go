package repository

import (
	"context"

	"github.com/lautistr/social-pet/models"
)

type Repository interface {
	Close()
	InsertPost(ctx context.Context, post *models.Post) error
	ListPosts(ctx context.Context) ([]*models.Post, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}

func Close() {
	repository.Close()
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return repository.InsertPost(ctx, post)
}

func ListPosts(ctx context.Context) ([]*models.Post, error) {
	return repository.ListPosts(ctx)
}
