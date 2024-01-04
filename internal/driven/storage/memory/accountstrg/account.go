package accountstrg

import (
	"context"

	"github.com/bad33ndj3/bunq2ynab/internal/core/entity"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Storage struct {
	data []*entity.Account
}

func (s *Storage) GetAccountByName(_ context.Context, name string) (*entity.Account, error) {
	res := lo.Filter(s.data, func(a *entity.Account, _ int) bool {
		return a.Description == name
	})

	if len(res) == 0 {
		return nil, errors.New("account not found")
	}

	if len(res) > 1 {
		return nil, errors.New("multiple accounts found")
	}

	return res[0], nil
}

func (s *Storage) SaveAccount(_ context.Context, b entity.Account) error {
	present := lo.Contains(s.data, &b)
	if present {
		return nil
	}

	s.data = append(s.data, &b)
	return nil
}

func New() (*Storage, error) {
	return &Storage{data: make([]*entity.Account, 0)}, nil
}
