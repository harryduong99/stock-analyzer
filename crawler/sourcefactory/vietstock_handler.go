package sourcefactory

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const formatMDY = "01/02/2006"

type VietstockSourceHandler struct {
}

func (sourceHandler VietstockSourceHandler) GetData(stocks []string, totalDays int, driver string) {
	for _, stock := range stocks {
		getVietstock(stock, totalDays)
	}
}

func getVietstock(stock string, totalDays int) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var url = `https://finance.vietstock.vn/` + stock + `/transaction-statistics.htm`
	fmt.Println("Visiting", url)
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text(getVietstockNodeName(1), &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	infos := strings.Fields(res)
	storeVietstockData(stock, infos)
}

func storeVietstockData(stock string, infos []string) ([]string, error) {
	stockInfos := []models.StockInfo{}
	stockInfo := models.StockInfo{}

	stockInfo.Code = stock
	dt, _ := time.Parse(formatMDY, infos[0])
	stockInfo.Date = primitive.NewDateTimeFromTime(dt)
	stockInfo.Change = infos[5] + " (" + infos[6] + " %)"
	stockInfo.AdjustedPrice, _ = convertResultToFloat(infos[4])
	stockInfo.ClosedPrice = 0
	stockInfo.HighestPrice = 0
	stockInfo.OpenPrice = 0
	stockInfo.LowestPrice = 0
	stockInfo.StockOrderAmount, _ = convertResultToInt(infos[7])
	stockInfo.StockOrderValue, _ = convertResultToInt(infos[8])
	stockInfo.StockDealAmount = 0
	stockInfo.StockDealValue = 0

	stockInfos = append(stockInfos, stockInfo)

	storingErr := repository.StockRepository.StoreStocks(stockInfos)
	if storingErr == nil {
		fmt.Println("Crawled " + stock)
	}

	return infos, storingErr
}

func getVietstockNodeName(backDay int) string {
	return `#stock-trading-result table tbody tr:nth-child(` + strconv.Itoa(backDay) + `)`
}
