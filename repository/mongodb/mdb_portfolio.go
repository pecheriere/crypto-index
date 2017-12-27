package mongodb

import (
	"gopkg.in/mgo.v2"
	"github.com/pecheriere/crypto-index/models"
	"github.com/pecheriere/crypto-index/repository"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"time"
)

type mdbPortfolioRepository struct {
	session *mgo.Session
}

func NewMDBPortfolioRepository() (repository.PortfolioRepository, error) {
	session, err := GetSession()
	if err != nil {
		return &mdbPortfolioRepository{}, err
	}
	return &mdbPortfolioRepository{session}, nil
}

func (m *mdbPortfolioRepository) Insert(portfolio *models.Portfolio) (error) {
	c := m.session.DB("crypto-index").C("portfolios")
	return c.Insert(portfolio)
}

func (m *mdbPortfolioRepository) FindLast(name string) (models.Portfolio, error) {
	c := m.session.DB("crypto-index").C("portfolios")
	result := models.Portfolio{}
	err := c.Find(bson.M{"name": name}).Sort("-updated_at").Limit(1).One(&result)
	return result, err
}

func (m *mdbPortfolioRepository) GetByHours(name string, from time.Time, to time.Time) ([]models.Portfolio, error) {
	c := m.session.DB("crypto-index").C("portfolios")
	p := c.Pipe([]bson.M{
		{
			"$match": bson.M{
				"name": name,
				"updated_at": bson.M{
					"$gte": from,
					"$lte": to,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"hour":  bson.M{"$hour": "$updated_at"},
					"day":   bson.M{"$dayOfMonth": "$updated_at"},
					"month": bson.M{"$month": "$updated_at"},
					"year":  bson.M{"$year": "$updated_at"},
				},
				"avg": bson.M{"$avg": "$value_usd"},
			},
		},
	})
	var res []bson.M
	p.All(&res)
	fmt.Println(len(res), res)
	return nil, nil
}

func (m *mdbPortfolioRepository) GetAll(name string, from time.Time, to time.Time) ([]models.Portfolio, error) {
	c := m.session.DB("crypto-index").C("portfolios")
	var result []models.Portfolio
	err := c.Find(bson.M{
		"name": name,
		"updated_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}).Sort("-updated_at").All(&result)
	return result, err
}
