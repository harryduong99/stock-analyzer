package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/duongnam99/stock-analyzer/crawler"
	"github.com/duongnam99/stock-analyzer/databasedriver"
	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"github.com/joho/godotenv"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

func init() {
	loadEnv()
	databasedriver.Mongo.ConnectDatabase()
}

func main() {
	if getAction() == "crawl" {
		crawler.Crawl(getSource())
	}

	if getAction() == "analyze" {
		// for _, stock := range getStocks() {

		// }

		series := techan.NewTimeSeries()

		// fetch this from your preferred exchange
		dataset := [][]string{
			// Timestamp, Open, Close, High, Low, volume
			{"1234567", "1", "2", "3", "5", "6"},
		}

		for _, datum := range dataset {
			start, _ := strconv.ParseInt(datum[0], 10, 64)
			period := techan.NewTimePeriod(time.Unix(start, 0), time.Hour*24)

			candle := techan.NewCandle(period)
			candle.OpenPrice = big.NewFromString(datum[1])
			candle.ClosePrice = big.NewFromString(datum[2])
			candle.MaxPrice = big.NewFromString(datum[3])
			candle.MinPrice = big.NewFromString(datum[4])

			series.AddCandle(candle)
		}

		closePrices := techan.NewClosePriceIndicator(series)
		movingAverage := techan.NewEMAIndicator(closePrices, 10) // Create an exponential moving average with a window of 10

		fmt.Println(movingAverage.Calculate(0).FormattedString(2))
	}
}

func getStocks() []models.StockInfo {
	m := map[string]interface{}{"code": []string{"FPT"}, "date": "05/01/2022"}
	return repository.GetStocks(
		context.Background(),
		m,
		10,
	)
}

func getSource() string {
	return os.Args[2]
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
