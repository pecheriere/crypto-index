package models

import (
	"testing"
)

func TestSortPositionsByRank(t *testing.T) {
	var positions []Position

	positions = append(positions, Position{
		Rank:2,
	})

	positions = append(positions, Position{
		Rank:3,
	})

	positions = append(positions, Position{
		Rank:1,
	})

	portfolio := Portfolio{
		Positions: positions,
	}

	portfolio.SortPositionsByRank()

	if portfolio.Positions[0].Rank != 1 {
		t.Fail()
	}
}

func TestSellPositions(t *testing.T) {
	var positions []Position

	positions = append(positions, Position{
		Symbol: "POS1",
		Rank:1,
	})

	positions = append(positions, Position{
		Symbol: "POS2",
		ValueUSD: 5,
		Rank:2,
	})

	positions = append(positions, Position{
		Symbol: "POS3",
		Rank:3,
	})

	portfolio := Portfolio{
		Positions: positions,
		CashUSD: 2,
	}

	portfolio.SellPositions(portfolio.Positions[1])

	if portfolio.Positions[1].Rank != 3 {
		t.Fail()
	}
	if portfolio.CashUSD != 5 + 2 {
		t.Fail()
	}
}