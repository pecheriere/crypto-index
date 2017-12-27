package portfolio

import (
	"testing"
	"github.com/pecheriere/crypto-index/repository/mock"
	"github.com/pecheriere/crypto-index/api/mock"
	"github.com/pecheriere/crypto-index/models"
	"time"
	"github.com/pecheriere/crypto-index/api"
	"github.com/stretchr/testify/assert"
	"errors"
	"log"
)

func TestCreateEQPortfolio(t *testing.T) {
	repo, _ := mock_db.NewMockPortfolioRepository()
	apiCoin := mock_api.NewMockCoinAPI()
	CreateEQPortfolio("TestPortFolio1", 6, 1000, repo, apiCoin)
	portfolio, _ := repo.FindLast("TestPortFolio1")

	assert.Equal(t, portfolio.Name, "TestPortFolio1")
	assert.Len(t, portfolio.Positions, 6)
	assert.Equal(t, portfolio.CashUSD, 0.0)
}

func TestUpdateEQPortfolio(t *testing.T) {
	repo, _ := mock_db.NewMockPortfolioRepository()

	repo.Insert(&models.Portfolio{
		Name:      "TestPortFolio",
		CashUSD:   0.0,
		ValueUSD:  330,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Positions: []models.Position{
			models.Position{
				CostUSD:  100,
				ValueUSD: 110,
				Quantity: 15,
				Symbol:   "MOCK1",
				Name:     "MockCoin1",
				Rank:     1,
			},
			models.Position{
				CostUSD:  100,
				ValueUSD: 110,
				Quantity: 15,
				Symbol:   "MOCK2",
				Name:     "MockCoin2",
				Rank:     2,
			},
			models.Position{
				CostUSD:  100,
				ValueUSD: 110,
				Quantity: 15,
				Symbol:   "MOCK6",
				Name:     "MockCoin6",
				Rank:     3,
			},
		},
	})

	apiCoin := new(mock_api.MockCoinAPI)
	apiCoin.On("GetCoinData", "MOCK6").Return(api.Coin{
		ID:           "ID6",
		Name:         "MockCoin6",
		Symbol:       "MOCK%d6",
		PriceUsd:     114,
		MarketCapUsd: 1273761623,
		Rank:         6,
	}, nil)

	UpdateEQPortfolio("TestPortFolio", repo, apiCoin)

	portfolio, _ := repo.FindLast("TestPortFolio")

	portfolio.SortPositionsByRank()

	assert.Equal(t, portfolio.Name, "TestPortFolio")
	assert.Len(t, portfolio.Positions, 3)
	assert.Equal(t, portfolio.Positions[2].Rank, 3)






	apiCoin.On("GetCoinData", "MOCK6").Return(api.Coin{}, errors.New("not found"))

	UpdateEQPortfolio("TestPortFolio", repo, apiCoin)

	portfolio, _ = repo.FindLast("TestPortFolio")

	portfolio.SortPositionsByRank()

	log.Println(portfolio)

	assert.Equal(t, portfolio.Name, "TestPortFolio")
	assert.Len(t, portfolio.Positions, 3)
	assert.Equal(t, portfolio.Positions[2].Rank, 3)

	t.Fail()
}
