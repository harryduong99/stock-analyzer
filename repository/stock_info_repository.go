package repository

import (
	"context"
	"log"

	"github.com/duongnam99/stock-analyzer/config"
	databasedriver "github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StockRepo struct{}

var StockRepository = &StockRepo{}

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

	for _, stock := range stocks {
		if !StockRepository.IsStockInfoExisting(context.Background(), stock.Code, stock.Date) {
			bbytes, _ := bson.Marshal(stock)
			docs = append(docs, bbytes)
		}
	}
	_, err := collection.InsertMany(context.Background(), docs)

	if err != nil {
		return err
	}

	return nil
}

func (stockRepo *StockRepo) IsStockInfoExisting(ctx context.Context, code string, date primitive.DateTime) bool {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)
	var stock models.StockInfo
	data := collection.FindOne(ctx, bson.M{"code": code, "date": date})
	err := data.Decode(&stock)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (stockRepo *StockRepo) GetStock(ctx context.Context, params map[string]interface{}) (models.StockInfo, error) {
	var stock models.StockInfo
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)
	data := collection.FindOne(ctx, bson.M{
		"code": params["code"],
		"date": params["date"],
	})
	error := data.Decode(&stock)
	return stock, error
}

func (stockRepo *StockRepo) GetLastDayStock(ctx context.Context, code string) models.StockInfo {
	var stock models.StockInfo
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	opts := options.FindOne().SetSort(bson.M{"date": -1})
	if err := collection.FindOne(ctx, bson.M{"code": code}, opts).Decode(&stock); err != nil {
		log.Fatal(err)
	}

	return stock
}

func (stockRepo *StockRepo) GetAndSort(ctx context.Context, code string, limit int, isDesc bool) []models.StockInfo {
	var stock models.StockInfo
	var stocks []models.StockInfo
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	opts := options.Find().SetSort(bson.M{"date": -1}).SetLimit(int64(11))
	cur, err := collection.Find(ctx, bson.M{"code": code}, opts)

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

func (stockRepo *StockRepo) GetStockByTime(ctx context.Context, params map[string]interface{}) []models.StockInfo {
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

func (stockRepo *StockRepo) GetStockDates(ctx context.Context, params map[string]interface{}, limit int) []models.StockInfo {
	var stock models.StockInfo
	var stocks []models.StockInfo

	option := options.Find().SetSort(bson.M{"date": -1}).SetLimit(int64(limit))
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK)

	cur, err := collection.Find(ctx, bson.M{
		"code": params["code"],
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

func (stockRepo *StockRepo) GetStocks(ctx context.Context, params map[string]interface{}, limit int) []models.StockInfo {
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
