package mock_api

import (
	"github.com/pecheriere/crypto-index/api"
	"fmt"
	"math/rand"
	"github.com/stretchr/testify/mock"
)

type MockCoinAPI struct {
	mock.Mock
}

func NewMockCoinAPI() (api.CoinApi) {
	return &MockCoinAPI{}
}

func (m *MockCoinAPI) GetCoinData(coin string) (api.Coin, error) {
	args := m.Called(coin)
	return args.Get(0).(api.Coin), args.Error(1)
}

func (m *MockCoinAPI) GetCoinsByMarketCap(toRank int) ([]api.Coin, error) {
	var result []api.Coin
	for i := 0; i < toRank; i++ {
		result = append(result, api.Coin{
			ID:           fmt.Sprintf("ID%d", i),
			Name:         fmt.Sprintf("MockCoin%d", i+1),
			Symbol:       fmt.Sprintf("MOCK%d", i+1),
			PriceUsd:     float64(rand.Intn(101)),
			MarketCapUsd: float64(986532 - (i * 10)),
			Rank:         i + 1,
		})
	}
	return result, nil
}
