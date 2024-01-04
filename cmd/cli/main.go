package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/bad33ndj3/bunq2ynab/internal/core/service/sync"
	"github.com/bad33ndj3/bunq2ynab/internal/driven/bunq"
	"github.com/bad33ndj3/bunq2ynab/internal/driven/storage/memory/accountstrg"
	iynab "github.com/bad33ndj3/bunq2ynab/internal/driven/ynab"
	"github.com/bad33ndj3/bunq2ynab/internal/driver/cli"
	"github.com/brunomvsouza/ynab.go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("Successfully synced!")
}

func run() error {
	ctx := context.Background()
	cfg, err := setupConfig()
	if err != nil {
		return errors.Wrap(err, "setting up config")
	}

	sv, err := setupSyncService(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "setting up sync service")
	}

	err = cli.NewClient(sv).Sync(ctx, time.Date(2024, 1, 3, 21, 0, 0, 0, time.UTC))
	if err != nil {
		return errors.Wrap(err, "syncing")
	}

	return nil
}

func setupSyncService(ctx context.Context, cfg *entity.Config) (*sync.Client, error) {
	bq, err := bunq.NewClient(ctx, cfg.BunqToken)
	if err != nil {
		return nil, errors.Wrap(err, "creating bunq client")
	}

	yn := iynab.NewClient(ynab.NewClient(cfg.YnabToken))

	bqs, err := accountstrg.New()
	if err != nil {
		return nil, errors.Wrap(err, "creating account storage")
	}
	sv := sync.NewClient(bq, bqs, yn, cfg)

	return sv, nil
}

func setupConfig() (*entity.Config, error) {
	dat, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	cfg := &entity.Config{}
	err = yaml.Unmarshal(dat, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling config file")
	}
	return cfg, nil
}
