package mock_db

import (
	"github.com/pecheriere/crypto-index/repository"
	"github.com/pecheriere/crypto-index/models"
	"errors"
	"time"
)

type mockPortfolioRepository struct {
	db []models.Portfolio
}

func NewMockPortfolioRepository() (repository.PortfolioRepository, error) {
	return &mockPortfolioRepository{}, nil
}

func (m *mockPortfolioRepository) Insert(portfolio *models.Portfolio) (error) {
	m.db = append(m.db, *portfolio)
	return nil
}

func (m *mockPortfolioRepository) FindLast(name string) (models.Portfolio, error) {
	var result *models.Portfolio

	for _, portfolio := range m.db {
		if portfolio.Name == name {
			if result == nil || result.UpdatedAt.Before(portfolio.UpdatedAt){
				result = &portfolio
			}
		}
	}

	if result == nil {
		return models.Portfolio{}, errors.New("not found")
	}

	return *result, nil
}

func (m *mockPortfolioRepository) GetByHours(name string, from time.Time, to time.Time) ([]models.Portfolio, error) {
	return nil, errors.New("not implemented")
}
