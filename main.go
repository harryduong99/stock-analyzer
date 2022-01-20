package main

// * * * * * cd work/src/github.com/duongnam99/stock-analyzer && ./stock-analyzer analyze
// 0 0 * * 1,2,3,4,5 <user> <command>
// 10 16 * * 1-5 /path_to_script

import (
	"log"
	"os"
	"strconv"

	"github.com/duongnam99/stock-analyzer/analyzer"
	"github.com/duongnam99/stock-analyzer/config"
	"github.com/duongnam99/stock-analyzer/crawler"
	"github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/httpcore"
	"github.com/duongnam99/stock-analyzer/httpcore/mail"
	"github.com/joho/godotenv"
)

func init() {
	loadEnv()
	databasedriver.Mongo.ConnectDatabase()
}

func main() {
	if len(os.Args) == 1 {
		httpcore.InitRoutes()
	}

	if getAction() == "crawl" {
		crawler.Crawl(getSource(), getTotalDays())
	}

	if getAction() == "analyze" {
		totalDays, _ := strconv.Atoi(os.Args[2])
		crawler.Crawl(config.CAFEF, totalDays)
		crawler.Crawl(config.VIETSTOCK, totalDays)

		results := analyzer.Analyze()
		mail.SendAnalyzeResult(results)
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
