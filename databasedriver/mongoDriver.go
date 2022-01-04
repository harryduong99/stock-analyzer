package databasedriver

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client *mongo.Client
}

var Mongo = &MongoDB{}

func (mongodb *MongoDB) ConnectDatabase() {
	// connStr := getConnectionString(user, password)
	client, err := mongo.NewClient(options.Client().ApplyURI(getConnectionString()))

	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	fmt.Println("connection ok")
	Mongo.Client = client
}

func (mongodb *MongoDB) ConnectCollection(databaseName, collectionName string) *mongo.Collection {
	return mongodb.Client.Database(databaseName).Collection(collectionName)
}

func getConnectionString() string {

	if os.Getenv("LOCAL_MODE") == "on" {
		return os.Getenv("MONGODB_CONNECTION_LOCAL")
	}

	connStr := os.Getenv("MONGODB_CONNECTION_ONL")

	return connStr
}
