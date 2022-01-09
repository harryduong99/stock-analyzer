package analyzer

import (
	"context"
	"math"

	"github.com/duongnam99/stock-analyzer/repository"
)

func GetAverageAmount(code string, totalDays int) float64 {
	stocks := repository.StockRepository.GetStockDates(
		context.Background(),
		map[string]interface{}{"code": code},
		totalDays,
	)
	totalAmount := 0
	for _, stock := range stocks {
		totalAmount += stock.StockDealAmount + stock.StockOrderAmount
	}

	return math.Round(float64(totalAmount/totalDays)*100) / 100
}

func GetFluctuatedFromAverageAmount(code string, totalDays int) float64 {
	averageAmount := GetAverageAmount(code, totalDays)
	lastdayStock := repository.StockRepository.GetLastDayStock(
		context.Background(),
		code,
	)

	amount := lastdayStock.StockDealAmount + lastdayStock.StockOrderAmount

	return math.Round(((float64(amount)-averageAmount)/averageAmount)*100) / 100
}

func GetAveragePrice(code string, totalDays int) float64 {
	stocks := repository.StockRepository.GetStockDates(
		context.Background(),
		map[string]interface{}{"code": code},
		totalDays,
	)
	var totalPrice float64
	totalPrice = 0
	for _, stock := range stocks {
		totalPrice += stock.AdjustedPrice
	}

	return math.Round(float64(totalPrice/float64(totalDays))*100) / 100
}

func GetFluctuatedFromAveragePrice(code string, totalDays int) float64 {
	averagePrice := GetAveragePrice(code, totalDays)
	lastdayStock := repository.StockRepository.GetLastDayStock(
		context.Background(),
		code,
	)

	price := lastdayStock.AdjustedPrice

	return math.Round(((float64(price)-averagePrice)/averagePrice)*100) / 100
}
