package writer

import (
	"generator/internal/models"
	"github.com/go-redis/redis/v8"

	"context"
)

type Writer struct {
	rdb *redis.Client
}

func NewWriter(rdb *redis.Client) *Writer {
	return &Writer{rdb: rdb}
}

// Write sends messages in stream
func (g *Writer) Write(ctx context.Context, stock *models.Stock) error {
	err := g.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "stream",
		Values: []interface{}{"ID", stock.ID, "Title", stock.Title, "Price", stock.Price},
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
