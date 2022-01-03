package sourcefactory

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type CafefSourceHandler struct {
}

func (sourceHandler CafefSourceHandler) GetData(stocks []string, driver string) {
	for _, stock := range stocks {
		if driver == "chrome" {
			result, error := GetCafefByChrome(stock, 0)
			if error != nil {
				log.Printf("Can not get data of %s", stock)
			}
			fmt.Println(result)
		} else {
			get(stock)
		}

	}

}

func get(stock string) {
	stockInfos := []StockInfo{}
	url := "https://s.cafef.vn/Lich-su-giao-dich-" + stock + "-1.chn"

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#ctl00_ContentPlaceHolder1_ctl03_rptData2_ctl01_itemTR", func(e *colly.HTMLElement) {
		stockInfo := StockInfo{}
		// data := strings.Fields(e.ChildText("td"))

		e.ForEach("td", func(i int, el *colly.HTMLElement) {
			if i > 11 {
				return
			}

			switch i {
			case 0:
				stockInfo.Date = el.Text
			case 1:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.AdjustedPrice = value
				}
			case 2:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.ClosedPrice = value
				}
			case 3:
				stockInfo.Change = el.Text
			case 5:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.StockOrderAmount = value
				}
			case 6:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.StockOrderValue = value
				}
			case 7:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.StockDealAmount = value
				}
			case 8:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.StockDealValue = value
				}
			case 9:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.OpenPrice = value
				}

			case 10:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.HighestPrice = value
				}
			case 11:
				if value, err := convertResult(el.Text); err == nil {
					stockInfo.LowestPrice = value
				}
			}
		})
		stockInfos = append(stockInfos, stockInfo)
	})

	c.OnScraped(func(r *colly.Response) {
		data, err := json.Marshal(stockInfos)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Finished. Here is your data:", string(data))
		}

	})

	c.Visit(url)

}

func convertResult(s string) (float32, error) {
	trimedSpace := strings.TrimSpace(s)
	rs, err := strconv.ParseFloat(trimedSpace, 32)
	if rs != 0 {
		return float32(rs), err
	}
	rs, err = strconv.ParseFloat(strings.Replace(trimedSpace, ",", "", -1), 32)

	return float32(rs), err
}
