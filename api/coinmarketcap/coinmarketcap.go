package coinmarketcap

import (
	"fmt"
	"encoding/json"
	"github.com/pecheriere/crypto-index/api"
	"log"
	"errors"
)

type coinCoinMarketCap struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUsd         float64 `json:"price_usd,string"`
	PriceBtc         float64 `json:"price_btc,string"`
	Usd24hVolume     float64 `json:"24h_volume_usd,string"`
	MarketCapUsd     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1h  float64 `json:"percent_change_1h,string"`
	PercentChange24h float64 `json:"percent_change_24h,string"`
	PercentChange7d  float64 `json:"percent_change_7d,string"`
	LastUpdated      int     `json:"last_updated,string"`
}

type coinMaketCapAPI struct {
	baseUrl string
	coins   []api.Coin
}

func NewCoinMaketCapAPI() (api.CoinApi) {
	return &coinMaketCapAPI{
		baseUrl: "https://api.coinmarketcap.com/v1",
	}
}

func (c *coinMaketCapAPI) GetCoinData(name string) (api.Coin, error) {
	id, err := c.getIdForSymbol(name)
	if err != nil {
		return api.Coin{}, err
	}

	url := fmt.Sprintf("%s/ticker/%s", c.baseUrl, id)
	resp, err := makeReq(url)
	if err != nil {
		return api.Coin{}, err
	}
	var data []coinCoinMarketCap
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return api.Coin{}, err
	}

	return convertToCoin(&data[0]), nil
}

func (c *coinMaketCapAPI) GetCoinsByMarketCap(toRank int) ([]api.Coin, error) {
	url := fmt.Sprintf("%s/ticker/?limit=%d", c.baseUrl, toRank)
	resp, err := makeReq(url)
	if err != nil {
		return nil, err
	}
	var data []coinCoinMarketCap
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	var result []api.Coin

	for _, coin := range data {
		result = append(result, convertToCoin(&coin))
	}

	return result, nil
}

func convertToCoin(coinCoinMarketCap *coinCoinMarketCap) (api.Coin) {
	return api.Coin{
		ID:           coinCoinMarketCap.ID,
		Name:         coinCoinMarketCap.Name,
		Symbol:       coinCoinMarketCap.Symbol,
		PriceUsd:     coinCoinMarketCap.PriceUsd,
		MarketCapUsd: coinCoinMarketCap.MarketCapUsd,
		Rank:         coinCoinMarketCap.Rank,
	}
}

func (c *coinMaketCapAPI) getIdForSymbol(symbol string) (string, error) {
	//Caching
	if c.coins == nil {
		log.Println("Caching coins from CoinMarketAPI...")
		coins, err := c.GetCoinsByMarketCap(0)
		if err != nil {
			return "", err
		}
		c.coins = coins
	}

	log.Println("Searching", symbol)

	//Search
	for _, coin := range c.coins {
		if coin.Symbol == symbol {
			log.Println("Found", symbol)
			return coin.ID, nil
		}
	}

	//Not found, Re-Caching
	coins, err := c.GetCoinsByMarketCap(0)
	if err != nil {
		return "", err
	}
	c.coins = coins

	for _, coin := range c.coins {
		if coin.Symbol == symbol {
			return coin.ID, nil
		}
	}

	//Not found
	return "", errors.New("not found")
}
