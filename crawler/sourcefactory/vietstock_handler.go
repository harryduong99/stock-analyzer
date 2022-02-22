package sourcefactory

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const formatMDY = "01/02/2006"

var (
	vietstockWaitGroup = sync.WaitGroup{}
)

type VietstockSourceHandler struct {
}

func (sourceHandler VietstockSourceHandler) GetData(stocks []string, totalDays int, driver string) {
	vietstockWaitGroup.Add(len(stocks))
	for _, stock := range stocks {
		go getVietstock(stock, totalDays)
	}
	vietstockWaitGroup.Wait()
}

func getVietstock(stock string, totalDays int) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug using
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//Initialization parameters, first pass an empty data
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	//Execute an empty task to create a chrome instance in advance
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//Create a context with a timeout of 40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var url = `https://finance.vietstock.vn/` + stock + `/transaction-statistics.htm`
	fmt.Println("Visiting", url)
	var res string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.Text(getVietstockNodeName(1), &res, chromedp.NodeVisible),
	)
	if err != nil {
		return
	}

	infos := strings.Fields(res)
	storeVietstockData(stock, infos)
	vietstockWaitGroup.Done()
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
