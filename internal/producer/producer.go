// Package producer sends data in repository
package producer

import (
	"github.com/chucky-1/generator/internal/model"
	"github.com/chucky-1/generator/internal/repository"

	"context"
	"math/rand"
	"time"
)

const (
	ctxForWrite = time.Second
	min         = 99
	max         = 102
)

// Producer contains struct of repository
type Producer struct {
	rep *repository.Repository
}

// NewProducer is repository
func NewProducer(rep *repository.Repository) *Producer {
	return &Producer{rep: rep}
}

// Write executes business logic and calls repository
func (p *Producer) Write(symbol *model.Symbol) error {
	// update price of symbol
	rate := float32(rand.Intn(max-min)+min) / 100
	symbol.Bid *= rate
	symbol.Ask = symbol.Bid - (symbol.Bid * 0.02)

	ctx, cancel := context.WithTimeout(context.Background(), ctxForWrite)
	defer cancel()
	return p.rep.Write(ctx, symbol)
}
