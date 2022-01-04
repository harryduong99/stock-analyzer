package main

import (
	"log"
	"os"

	"github.com/duongnam99/stock-analyzer/crawler"
	"github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	databasedriver.Mongo.ConnectDatabase()
	crawler.Crawl(getSource())
}

func getSource() string {
	arg := os.Args[1:2]
	source := arg[0]
	return source
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
