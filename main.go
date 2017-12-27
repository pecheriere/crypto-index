package main

import (
	"github.com/pecheriere/crypto-index/engine/portfolio"
	"github.com/pecheriere/crypto-index/repository/mongodb"
	"github.com/pecheriere/crypto-index/api/coinmarketcap"
	"log"
)

func main() {
	portfolioRepo, err := mongodb.NewMDBPortfolioRepository()
	if err != nil {
		panic(err)
	}

	api := coinmarketcap.NewCoinMaketCapAPI()

	//portfolio.CreateEQPortfolio("PE10EQ", 10, 1000, portfolioRepo, api)
	//portfolio.CreateEQPortfolio("PE25EQ", 25, 1000, portfolioRepo, api)

	log.Println("Starting Data Fetching...")
	portfolio.LaunchDataFetching(portfolioRepo, api)

	//r := mux.NewRouter()
	//r.HandleFunc("/PE10EQ", GetPE10EQ)
	//log.Fatal(http.ListenAndServe(":8080", r))
}