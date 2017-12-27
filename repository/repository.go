package repository

import (
	"github.com/pecheriere/crypto-index/models"
	"time"
)

type PortfolioRepository interface {
	Insert(portfolio *models.Portfolio) (error)
	FindLast(name string) (models.Portfolio, error)
	GetByHours(name string, from time.Time, to time.Time) ([]models.Portfolio, error)
}
