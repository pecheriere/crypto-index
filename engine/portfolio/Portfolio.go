package portfolio

import (
	"github.com/pecheriere/crypto-index/models"
	"time"
	"github.com/pecheriere/crypto-index/repository"
	"github.com/pecheriere/crypto-index/api"
	"log"
)

func CreateEQPortfolio(name string, ranks int, cash float64, repo repository.PortfolioRepository, api api.CoinApi) {
	if coins, err := api.GetCoinsByMarketCap(ranks); err != nil {
		log.Println(err)
	} else {
		budgetPerCoin := cash / float64(len(coins))
		portfolio := models.Portfolio{}
		portfolio.Name = name
		portfolio.ValueUSD = cash
		portfolio.CashUSD = 0
		portfolio.CreatedAt = time.Now().UTC()
		portfolio.UpdatedAt = portfolio.CreatedAt
		for _, coin := range coins {
			portfolio.Positions = append(portfolio.Positions, models.Position{
				CostUSD:  budgetPerCoin,
				ValueUSD: budgetPerCoin,
				Quantity: budgetPerCoin / coin.PriceUsd,
				Symbol:   coin.Symbol,
				Name:     coin.Name,
				Rank:     coin.Rank,
			})
		}
		repo.Insert(&portfolio)
	}
}

func updateEQPortfolioWithCoins(name string, repo repository.PortfolioRepository, api api.CoinApi, coins []api.Coin) {
	if portfolio, err := getUpdatedPortfolioWithCoins(name, repo, api, coins); err != nil {
		log.Println("updateEQPortfolioWithCoins", err)
	} else {
		portfolioLen := len(portfolio.Positions)
		for i := len(portfolio.Positions) - 1; i >= 0; i-- {
			if portfolio.Positions[i].Rank > portfolioLen {
				portfolio.SellPositions(portfolio.Positions[i])
			} else {
				break
			}
		}

		budgetPerCoin := portfolio.CashUSD / float64(portfolioLen-len(portfolio.Positions)) // Number of missing positions
		portfolio.ValueUSD = 0

		for i := 0; i < portfolioLen; i++ {
			found := false
			for _, position := range portfolio.Positions {
				if position.Symbol == coins[i].Symbol {
					found = true
					portfolio.ValueUSD += position.ValueUSD
					break
				}
			}
			if !found {
				portfolio.Positions = append(portfolio.Positions, models.Position{
					CostUSD:  budgetPerCoin,
					ValueUSD: budgetPerCoin,
					Quantity: budgetPerCoin / coins[i].PriceUsd,
					Symbol:   coins[i].Symbol,
					Name:     coins[i].Name,
					Rank:     coins[i].Rank,
				})
				portfolio.ValueUSD += budgetPerCoin
			}
		}
		portfolio.CashUSD = 0
		repo.Insert(&portfolio)
	}
}

func UpdateEQPortfolio(name string, repo repository.PortfolioRepository, api api.CoinApi) {
	if portfolio, err := getUpdatedPortfolio(name, repo, api); err != nil {
		log.Println("UpdateEQPortfolio", err)
	} else {
		portfolioLen := len(portfolio.Positions)
		if coins, err := api.GetCoinsByMarketCap(portfolioLen); err != nil {
			log.Println(err)
		} else {
			updateEQPortfolioWithCoins(name, repo, api, coins)
		}
	}
}

func getUpdatedPortfolioWithCoins(name string, repo repository.PortfolioRepository, api api.CoinApi, coins []api.Coin) (models.Portfolio, error) {
	portfolio, err := repo.FindLast(name);
	if err != nil {
		return portfolio, err
	}

	portfolio.ValueUSD = 0
	portfolio.UpdatedAt = time.Now().UTC()
	for idx := range portfolio.Positions {
		found := false
		for _, coin := range coins {
			if coin.Name == portfolio.Positions[idx].Name {
				found = true
				portfolio.Positions[idx].ValueUSD = portfolio.Positions[idx].Quantity * coin.PriceUsd
				portfolio.Positions[idx].Rank = coin.Rank
				portfolio.ValueUSD += portfolio.Positions[idx].ValueUSD
				break;
			}
		}
		if !found {
			if coin, err := api.GetCoinData(portfolio.Positions[idx].Symbol); err != nil {
				log.Println("getUpdatedPortfolioWithCoins",
					"Error while looking for", portfolio.Positions[idx].Name, err)
			} else {
				portfolio.Positions[idx].ValueUSD = portfolio.Positions[idx].Quantity * coin.PriceUsd
				portfolio.Positions[idx].Rank = coin.Rank
				portfolio.ValueUSD += portfolio.Positions[idx].ValueUSD
			}
		}
	}

	portfolio.SortPositionsByRank()
	return portfolio, nil
}

func getUpdatedPortfolio(name string, repo repository.PortfolioRepository, api api.CoinApi) (models.Portfolio, error) {
	portfolio, err := repo.FindLast(name);
	if err != nil {
		return portfolio, err
	}

	coins, err := api.GetCoinsByMarketCap(len(portfolio.Positions))
	if err != nil {
		return portfolio, err
	}

	return getUpdatedPortfolioWithCoins(name, repo, api, coins)
}

func UpdatePortfolioValue(name string, repo repository.PortfolioRepository, api api.CoinApi) {
	if portfolio, err := getUpdatedPortfolio(name, repo, api); err != nil {
		log.Println(err)
	} else {
		repo.Insert(&portfolio)
	}
}

func GetPortfolioLastValue(name string, repo repository.PortfolioRepository) (models.Portfolio, error) {
	if portfolio, err := repo.FindLast(name); err != nil {
		return models.Portfolio{}, err
	} else {
		return portfolio, nil
	}
}

func GetLast24Hours(name string, repo repository.PortfolioRepository) ([]models.Portfolio, error) {
	from := time.Now().UTC().Add(time.Hour * -24)
	to := time.Now().UTC()
	if portfolio, err := repo.GetByHours("PE10EQ", from, to); err != nil {
		return nil, err
	} else {
		return portfolio, nil
	}
}

func LaunchDataFetching(repo repository.PortfolioRepository, api api.CoinApi) {
	for {
		if coins, err := api.GetCoinsByMarketCap(30); err != nil {
			log.Fatal(err)
		} else {
			updateEQPortfolioWithCoins("PE10EQ", repo, api, coins)
			updateEQPortfolioWithCoins("PE25EQ", repo, api, coins)
		}
		log.Println("Portfolios updated")
		<-time.After(6 * time.Minute)
	}
}
