package repository

import (
	"context"
	"log"

	"github.com/duongnam99/stock-analyzer/config"
	databasedriver "github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetStock(ctx context.Context, params map[string]interface{}) (models.StockInfo, error) {
	var stock models.StockInfo
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)
	data := collection.FindOne(ctx, bson.M{
		"code": params["code"],
		"date": params["date"],
	})
	error := data.Decode(&stock)
	return stock, error
}

func GetStockByTime(ctx context.Context, params map[string]interface{}) []models.StockInfo {
	var stock models.StockInfo
	var stocks []models.StockInfo

	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	cur, err := collection.Find(ctx, bson.M{
		"code": params["code"],
		"date": bson.M{"$in": params["date"]},
	})

	if err != nil {
		log.Println(err)
		return nil
	}
	for cur.Next(ctx) {
		cur.Decode(&stock)
		stocks = append(stocks, stock)
	}
	return stocks
}

func GetStocks(ctx context.Context, params map[string]interface{}, limit int) []models.StockInfo {
	var stock models.StockInfo
	var stocks []models.StockInfo
	option := options.Find().SetLimit(int64(limit))
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)
	cur, err := collection.Find(ctx, bson.M{
		"code": bson.M{"$in": params["code"]},
		"date": params["date"],
	}, option)

	if err != nil {
		log.Println(err)
		return nil
	}
	for cur.Next(ctx) {
		cur.Decode(&stock)
		stocks = append(stocks, stock)
	}
	return stocks
}
