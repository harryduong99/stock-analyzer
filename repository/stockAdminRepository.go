package repository

import (
	"context"
	"log"

	"github.com/duongnam99/stock-analyzer/config"
	databasedriver "github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/models"
	"go.mongodb.org/mongo-driver/bson"
)

type StockAdminRepo struct{}

var StockAdminRepository = &StockAdminRepo{}

func (stockRepo *StockAdminRepo) AllStockAdmin() []models.StockAdmin {
	var stocks []models.StockAdmin
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK_ADMIN)

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var stock models.StockAdmin
		// & character returns the memory address of the following variable.
		err := cur.Decode(&stock) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		stocks = append(stocks, stock)
	}

	return stocks
}

func (stockRepo *StockAdminRepo) StoreStockAdmin(stock models.StockAdmin) error {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK_ADMIN)

	bbytes, _ := bson.Marshal(stock)
	_, err := collection.InsertOne(context.Background(), bbytes)

	if err != nil {
		return err
	}

	return nil
}
