package main

import (
	"net/http"
	"github.com/pecheriere/crypto-index/repository/mongodb"
	"github.com/pecheriere/crypto-index/engine/portfolio"
	"encoding/json"
)

func GetPE10EQ(w http.ResponseWriter, r *http.Request) {
	portfolioRepo, _ := mongodb.NewMDBPortfolioRepository()
	p, _ := portfolio.GetPortfolioLastValue("PE10EQ", portfolioRepo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}
