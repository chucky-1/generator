package main

import (
	"generator/internal/configs"
	"generator/internal/models"
	"generator/internal/writer"

	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"context"
	"math/rand"
	"strconv"
	"time"
)

const (
	countOfStocks   = 10
	maxPriceOfStock = 1000
	min             = 99
	max             = 102
	updateTime      = time.Second / 2
	ctxForWrite     = time.Second
)

func write(w *writer.Writer, stock *models.Stock) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxForWrite)
	defer cancel()
	err := w.Write(ctx, stock)
	if err != nil {
		return err
	}
	return nil
}

// generate updates price of stock
func generate(stock *models.Stock) {
	rate := float32(rand.Intn(max - min) + min) / 100
	stock.Price *= rate
}

func main() {
	// Configuration
	cfg := new(configs.Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	// Redis
	hostAndPort := cfg.Host + ":" + cfg.Port
	rdb := redis.NewClient(&redis.Options{Addr: hostAndPort})

	w := writer.NewWriter(rdb)

	// Initial stocks
	var stocks []*models.Stock
	for i := 1; i < countOfStocks + 1; i++ {
		title := "stock " + strconv.Itoa(i)
		price := float32(rand.Intn(maxPriceOfStock))
		stock := models.Stock{
			ID:    i,
			Title: title,
			Price: price,
		}
		stocks = append(stocks, &stock)
	}

	// Business logic
	for {
		for _, stock := range stocks {
			generate(stock)
			err := write(w, stock)
			if err != nil {
				log.Error(err)
			}
			log.Infof("%s costs %f", stock.Title, stock.Price)
		}
		time.Sleep(updateTime)
	}
}
