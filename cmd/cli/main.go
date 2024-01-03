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

	log.Println(cfg)

	bq, err := bunq.NewClient(ctx, "key")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	yn := iynab.NewClient(ynab.NewClient("key"))

	err = cli.NewClient(bq, yn, cfg).Sync()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
