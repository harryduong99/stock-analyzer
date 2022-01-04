package repository

import (
	"context"

	"github.com/duongnam99/stock-analyzer/config"
	databasedriver "github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/models"
	"go.mongodb.org/mongo-driver/bson"
)

type StockRepo struct{}

var StockRepoistory = &StockRepo{}

func (stockRepo *StockRepo) StoreStock(stock models.StockInfo) error {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	bbytes, _ := bson.Marshal(stock)
	_, err := collection.InsertOne(context.Background(), bbytes)

	if err != nil {
		return err
	}

	return nil
}

func (stockRepo *StockRepo) StoreStocks(stocks []models.StockInfo) error {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	docs := []interface{}{}

	for _, link := range stocks {
		bbytes, _ := bson.Marshal(link)
		docs = append(docs, bbytes)
	}

	_, err := collection.InsertMany(context.Background(), docs)

	if err != nil {
		return err
	}

	return nil
}
