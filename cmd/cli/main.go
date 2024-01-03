package main

import (
	"bunqtoynab/core/domain"
	"bunqtoynab/internal/driven/bunq"
	iynab "bunqtoynab/internal/driven/ynab"
	"bunqtoynab/internal/driver/cli"
	"context"
	"github.com/brunomvsouza/ynab.go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	dat, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	cfg := &domain.Config{}
	err = yaml.Unmarshal(dat, cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	bq, err := bunq.NewClient(ctx, cfg.BunqToken)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	yn := iynab.NewClient(ynab.NewClient(cfg.YnabToken))

	err = cli.NewClient(bq, yn, cfg).Sync(time.Date(2024, 1, 3, 18, 0, 0, 0, time.UTC))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
