package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/chucky-1/generator/internal/config"
	"github.com/chucky-1/generator/internal/model"
	"github.com/chucky-1/generator/internal/producer"
	"github.com/chucky-1/generator/internal/repository"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	countOfStocks   = 10
	maxPriceOfStock = 1000
	updateTime      = time.Second / 2
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

	// Initial stocks
	var stocks []*model.Stock
	for i := 0; i < countOfStocks; i++ {
		title := fmt.Sprint("stock", " ", strconv.Itoa(i+1))
		price := float32(rand.Intn(maxPriceOfStock))
		stock := model.Stock{
			ID:    i + 1,
			Title: title,
			Price: price,
		}
		stocks = append(stocks, &stock)
	}

	// Business logic
	t := time.NewTicker(updateTime)
	for {
		select {
		case <-t.C:
			for _, stock := range stocks {
				err := prod.Write(stock)
				if err != nil {
					log.Error(err)
				} else {
					log.Infof("%s costs %f", stock.Title, stock.Price)
				}
			}
		}
	}
}
