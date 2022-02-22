package sourcefactory

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const formatDMY = "02/01/2006"

var (
	cafefWaitGroup = sync.WaitGroup{}
)

type CafefSourceHandler struct {
}

func (sourceHandler CafefSourceHandler) GetData(stocks []string, totalDays int, driver string) {
	cafefWaitGroup.Add(len(stocks))
	for _, stock := range stocks {
		if driver == "chrome" {
			result, error := GetCafefByChrome(stock, 0)
			if error != nil {
				log.Printf("Can not get data of %s", stock)
			}
			fmt.Println(result)
		} else {
			getCafef(stock, totalDays)
		}
	}
	cafefWaitGroup.Wait()
}

func getCafef(stock string, totalDays int) {
	stockInfos := []models.StockInfo{}
	url := "https://s.cafef.vn/Lich-su-giao-dich-" + stock + "-1.chn"

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#GirdTable2", func(e *colly.HTMLElement) {
		stockInfo := models.StockInfo{}
		stockInfo.Code = stock
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			if i > totalDays+1 {
				return
			}

			if i > 1 {
				el.ForEach("td", func(i int, ef *colly.HTMLElement) {
					if i > 11 {
						return
					}

					switch i {
					case 0:
						// stockInfo.Date = ef.Text
						// preParsing := strings.Replace(ef.Text, "-", "/", -1)
						dt, _ := time.Parse(formatDMY, ef.Text)
						stockInfo.Date = primitive.NewDateTimeFromTime(dt)
					case 1:
						if value, err := convertResultToFloat(ef.Text); err == nil {
							stockInfo.AdjustedPrice = value
						}
					case 2:
						if value, err := convertResultToFloat(ef.Text); err == nil {
							stockInfo.ClosedPrice = value
						}
					case 3:
						stockInfo.Change = ef.Text
					case 5:
						if value, err := convertResultToInt(ef.Text); err == nil {
							stockInfo.StockOrderAmount = value
						}
					case 6:
						if value, err := convertResultToInt(ef.Text); err == nil {
							stockInfo.StockOrderValue = value
						}
					case 7:
						if value, err := convertResultToInt(ef.Text); err == nil {
							stockInfo.StockDealAmount = value
						}
					case 8:
						if value, err := convertResultToInt(ef.Text); err == nil {
							stockInfo.StockDealValue = value
						}
					case 9:
						if value, err := convertResultToFloat(ef.Text); err == nil {
							stockInfo.OpenPrice = value
						}

					case 10:
						if value, err := convertResultToFloat(ef.Text); err == nil {
							stockInfo.HighestPrice = value
						}
					case 11:
						if value, err := convertResultToFloat(ef.Text); err == nil {
							stockInfo.LowestPrice = value
						}
					}
				})
				stockInfos = append(stockInfos, stockInfo)
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		_, err := json.Marshal(stockInfos)
		if err != nil {
			fmt.Println(err)
		} else {
			storingErr := repository.StockRepository.StoreStocks(stockInfos)
			if storingErr == nil {
				fmt.Println("Crawled " + stock)
			}
			cafefWaitGroup.Done()
		}

	})

	c.Visit(url)

}

func convertResultToFloat(s string) (float64, error) {
	trimedSpace := strings.TrimSpace(s)
	rs, err := strconv.ParseFloat(trimedSpace, 32)
	return math.Round(float64(rs)*100) / 100, err
}

func convertResultToInt(s string) (int, error) {
	trimedSpace := strings.TrimSpace(s)

	rs, err := strconv.Atoi(strings.Replace(trimedSpace, ",", "", -1))

	return rs, err
}
