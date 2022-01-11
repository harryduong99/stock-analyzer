package admin

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/duongnam99/stock-analyzer/analyzer"
)

func Report(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	results := analyzer.Analyze()
	rootPath, _ := os.Getwd()
	var dataToMail [][]string
	var data []string
	for _, stock := range results {
		data = []string{}
		data = append(data, stock.Code)
		data = append(data, strconv.FormatFloat(stock.OpenPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.AdjustClosedPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.YesterdayAdjustPrice, 'f', -1, 64))
		data = append(data, strconv.FormatFloat(stock.DayBeforeYesterdayAdjustPrice, 'f', -1, 64))
		data = append(data, stock.Change)
		data = append(data, strconv.FormatFloat(math.Round(stock.FluctuatedAmount*100), 'f', -1, 64)+"%")
		data = append(data, strconv.FormatFloat(math.Round(stock.FluctuatedPrice*100), 'f', -1, 64)+"%")
		dataToMail = append(dataToMail, data)
	}
	body := ParseTemplate(rootPath+"/httpcore/admin/views/analyze_result.html", map[string]interface{}{
		"stocks":     dataToMail,
		"report_url": os.Getenv("APP_URL") + "/stock/report",
	})

	fmt.Fprint(w, body)
}

func ParseTemplate(fileName string, data interface{}) string {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		log.Fatalln(err)
	}
	return buffer.String()
}
