package models

import (
	"time"
	"sort"
)

type Position struct {
	CostUSD				float64	`bson:"cost_usd" json:"costUsd"`
	ValueUSD			float64	`bson:"value_usd" json:"valueUsd"`
	Quantity			float64	`bson:"quantity" json:"quantity"`
	Symbol				string	`bson:"symbol" json:"symbol"`
	Name				string	`bson:"name" json:"name"`
	Rank				int		`bson:"rank" json:"rank"`
}

type Portfolio struct {
	Name 		string		`bson:"name" json:"name"`
	CashUSD		float64		`bson:"cash_usd" json:"cashUsd"`
	ValueUSD	float64		`bson:"value_usd" json:"valueUsd"`
	Positions	[]Position	`bson:"positions" json:"positions"`
	UpdatedAt	time.Time	`bson:"updated_at" json:"updatedAt"`
	CreatedAt	time.Time	`bson:"created_at" json:"createdAt"`
}

func (p *Portfolio) SortPositionsByRank() {
	sort.Slice(p.Positions, func(i, j int) bool {
		return p.Positions[i].Rank < p.Positions[j].Rank
	})
}

func (p *Portfolio) SellPositions(position Position) {
	for idx := range p.Positions {
		if p.Positions[idx].Symbol == position.Symbol {
			p.CashUSD += p.Positions[idx].ValueUSD
			p.Positions = append(p.Positions[:idx], p.Positions[idx+1:]...) // Keep order
			break;
		}
	}
}
