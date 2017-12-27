package api

import "fmt"

type Coin struct {
	ID				string
	Name			string
	Symbol			string
	PriceUsd		float64
	MarketCapUsd	float64
	Rank			int
}

type CoinApi interface {
	GetCoinData(coin string) (Coin, error)
	GetCoinsByMarketCap(toRank int) ([]Coin, error)
}

func (c Coin) String() string {
	return fmt.Sprintf("Name: %s | Value: %f | Market Cap: %f | Rank: %d", c.Name, c.PriceUsd, c.MarketCapUsd, c.Rank)
}