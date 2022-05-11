package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/chucky-1/generator/internal/config"
	"github.com/chucky-1/generator/internal/model"
	"github.com/chucky-1/generator/internal/producer"
	"github.com/chucky-1/generator/internal/repository"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"fmt"
	"math/rand"
	"time"
)

const (
	countOfSymbols   = 5
	maxPriceOfSymbol = 1000
	updateTime      = time.Second * 3
)

func main() {
	// Configuration
	cfg := new(config.Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("%v", err)
	}

	// Redis
	hostAndPort := fmt.Sprint(cfg.Host, ":", cfg.Port)
	rdb := redis.NewClient(&redis.Options{Addr: hostAndPort})

	rep := repository.NewRepository(rdb)
	prod := producer.NewProducer(rep)

	// Initial symbols
	var symbols []*model.Symbol
	for i := 0; i < countOfSymbols; i++ {
		bid := float32(rand.Intn(maxPriceOfSymbol))
		symbol := model.Symbol{
			ID:  uuid.New(),
			Bid: bid,
			Ask: bid - (bid * 0.02),
		}
		symbols = append(symbols, &symbol)
	}

	// Business logic
	t := time.NewTicker(updateTime)
	for {
		select {
		case <-t.C:
			for _, symbol := range symbols {
				err := prod.Write(symbol)
				if err != nil {
					log.Error(err)
				} else {
					log.Infof("%d costs. Bid: %f, ask: %f", symbol.ID, symbol.Bid, symbol.Ask)
				}
			}
		}
	}
}
