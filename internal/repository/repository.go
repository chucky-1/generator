// Package repository sends messages about prices change into stream
package repository

import (
	"github.com/chucky-1/generator/internal/model"
	"github.com/go-redis/redis/v8"

	"context"
)

// Repository contains client of redis
type Repository struct {
	rdb *redis.Client
}

// NewRepository is constructor
func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

// Write sends messages in stream
func (g *Repository) Write(ctx context.Context, stock *model.Stock) error {
	return g.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: []interface{}{"ID", stock.ID, "Title", stock.Title, "Price", stock.Price},
	}).Err()
}
