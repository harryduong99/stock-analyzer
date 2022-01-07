package crawler

import (
	"github.com/duongnam99/stock-analyzer/crawler/sourcefactory"
	"github.com/duongnam99/stock-analyzer/repository"
)

func Crawl(source string, totalDays int) {
	targets := getTargets()
	sourcefactory.GetSourceHandlerFactory(source).GetData(targets, totalDays, "")
}

func getTargets() []string {
	stocks := repository.StockAdminRepository.AllStockAdmin()
	var result []string
	for _, stock := range stocks {
		result = append(result, stock.Code)
	}
	return result
}
