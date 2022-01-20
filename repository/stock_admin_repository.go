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
	if !StockAdminRepository.IsStockAdminExisting(context.Background(), stock.Code) {
		_, err := collection.InsertOne(context.Background(), bbytes)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func (stockRepo *StockAdminRepo) DeleteOneById(id primitive.ObjectID) bool {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK_ADMIN)

	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.TODO(), filter)
	return err == nil
}

func (stockRepo *StockAdminRepo) DeleteOneByCode(code string) bool {
	var stock models.StockAdmin
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK_ADMIN)
	if err := collection.FindOne(context.Background(), bson.M{"code": code}).Decode(&stock); err != nil {
		log.Fatal(err)
	}

	result := stockRepo.DeleteOneById(stock.ID)

	return result
}

func (stockRepo *StockAdminRepo) IsStockAdminExisting(ctx context.Context, code string) bool {
	collection := databasedriver.Mongo.ConnectCollection(config.DB_NAME, config.COL_STOCK_ADMIN)
	var stock models.StockInfo
	data := collection.FindOne(ctx, bson.M{"code": code})
	err := data.Decode(&stock)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}
