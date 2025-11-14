package transactions

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Transaction struct {
	Date        time.Time
	Description string
	Out         *money.Money
	In          *money.Money
}

type Service interface {
	ListTransactions() ([]Transaction, error)
}

type Repository interface {
	GetAllTransactions() ([]Transaction, error)
}

type service struct {
	repo Repository
}

func (s *service) ListTransactions() ([]Transaction, error) {
	return s.repo.GetAllTransactions()
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
