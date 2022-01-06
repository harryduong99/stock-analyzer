package main

import (
	"log"
	"os"
	"strconv"

	"github.com/duongnam99/stock-analyzer/analyzer"
	"github.com/duongnam99/stock-analyzer/crawler"
	"github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/joho/godotenv"
)

func init() {
	loadEnv()
	databasedriver.Mongo.ConnectDatabase()
}

func main() {
	if getAction() == "crawl" {
		crawler.Crawl(getSource(), getTotalDays())
	}

	if getAction() == "analyze" {
		analyzer.Analyze()
	}
}

func getSource() string {
	return os.Args[2]
}

func getTotalDays() int {
	rs, _ := strconv.Atoi(os.Args[3])
	return rs
}

func getAction() string {
	return os.Args[1]
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
