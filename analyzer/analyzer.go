package analyzer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"github.com/sdcoffey/big"
	"github.com/sdcoffey/techan"
)

type AnalyzeResult struct {
	Code                          string
	OpenPrice                     float64
	AdjustClosedPrice             float64
	YesterdayAdjustPrice          float64
	DayBeforeYesterdayAdjustPrice float64
	Change                        string
	FluctuatedAmount              float64
	FluctuatedPrice               float64
}

var (
	waitGroup     = sync.WaitGroup{}
	lock          = sync.Mutex{}
	analyzeResult = []AnalyzeResult{}
)

func Analyze() []AnalyzeResult {
	stocks := repository.StockAdminRepository.AllStockAdmin()
	waitGroup.Add(len(stocks))
	for _, stock := range stocks {
		go ProcessStock(stock) // problem : may cause too many goroutine!
	}
	waitGroup.Wait()
	return analyzeResult
}

func ProcessStock(stock models.StockAdmin) {
	stockInfos := repository.StockRepository.GetAndSort(context.Background(), stock.Code, 3, true)
	if len(stockInfos) == 0 {
		waitGroup.Done()
		return
	}
	result := AnalyzeResult{}
	result.Code = stockInfos[0].Code
	result.OpenPrice = stockInfos[0].OpenPrice
	result.AdjustClosedPrice = stockInfos[0].AdjustedPrice
	result.YesterdayAdjustPrice = stockInfos[1].AdjustedPrice
	result.DayBeforeYesterdayAdjustPrice = stockInfos[2].AdjustedPrice
	result.Change = stockInfos[0].Change
	result.FluctuatedAmount = GetFluctuatedFromAverageAmount(stock.Code, 10)
	result.FluctuatedPrice = GetFluctuatedFromAveragePrice(stock.Code, 10)
	lock.Lock()
	analyzeResult = append(analyzeResult, result)
	lock.Unlock()
	waitGroup.Done()
}

func getStockByTime(stock string, totalDays int) []models.StockInfo {
	return repository.StockRepository.GetStockDates(
		context.Background(),
		map[string]interface{}{"code": stock},
		totalDays,
	)
}

func Temp() {
	series := techan.NewTimeSeries()
	stocks := getStockByTime("HPG", 10)

	var dataset [][]string
	// layout := "02/01/2006"

	set := []string{}
	for _, stock := range stocks {
		// t, err := time.Parse(layout, stock.Date)
		// if err != nil {
		// 	log.Fatalln("Error while parsing date :", err)
		// }
		// set = append(set, strconv.Itoa(int(t.Unix())))
		set = append(set, fmt.Sprintf("%f", stock.OpenPrice))
		set = append(set, fmt.Sprintf("%f", stock.ClosedPrice))
		set = append(set, fmt.Sprintf("%f", stock.HighestPrice))
		set = append(set, fmt.Sprintf("%f", stock.LowestPrice))
		set = append(set, strconv.Itoa(stock.StockOrderAmount+stock.StockDealAmount))
		dataset = append(dataset, set)
		set = []string{}
	}

	// log.Fatalln(dataset)

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

	indicator := techan.NewMACDIndicator(techan.NewClosePriceIndicator(series), 12, 26)

	// log.Fatalln(macd)
	// assert.NotNil(t, macd)

	// closePrices := techan.NewClosePriceIndicator(series)
	// movingAverage := techan.NewEMAIndicator(closePrices, 10) // Create an exponential moving average with a window of 10

	// log.Fatalln(movingAverage.Calculate(0).FormattedString(2))

	// indicator := techan.NewClosePriceIndicator(series)

	// record trades on this object
	record := techan.NewTradingRecord()

	entryConstant := techan.NewConstantIndicator(30)
	exitConstant := techan.NewConstantIndicator(10)

	// Is satisfied when the price ema moves above 30 and the current position is new
	entryRule := techan.And(
		techan.NewCrossUpIndicatorRule(entryConstant, indicator),
		techan.PositionNewRule{})

	// Is satisfied when the price ema moves below 10 and the current position is open
	exitRule := techan.And(
		techan.NewCrossDownIndicatorRule(indicator, exitConstant),
		techan.PositionOpenRule{})

	strategy := techan.RuleStrategy{
		UnstablePeriod: 10, // Period before which ShouldEnter and ShouldExit will always return false
		EntryRule:      entryRule,
		ExitRule:       exitRule,
	}

	result := strategy.ShouldEnter(0, record) // returns false
	log.Fatalln(result)
}
