package sourcefactory

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"
)

func GetCafefByChrome(stock string, totalDays int) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var url = `https://s.cafef.vn/Lich-su-giao-dich-` + stock + `-1.chn`
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Text(getNodeName(totalDays), &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	infos := strings.Fields(res)

	return infos, err
}

func getNodeName(totalDays int) string {
	sub := strconv.Itoa(totalDays + 1)
	if totalDays < 9 {
		sub = "0" + strconv.Itoa(totalDays+1)
	}
	node := `#ctl00_ContentPlaceHolder1_ctl03_rptData2_ctl` + sub + `_itemTR`

	return node
}

// func Chrome() {
// opts := append(chromedp.DefaultExecAllocatorOptions[:],
// 	chromedp.Flag("headless", false),
// 	chromedp.Flag("disable-gpu", false),
// 	chromedp.Flag("enable-automation", false),
// 	chromedp.Flag("disable-extensions", false),
// )

// allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
// defer cancel()

// // create context
// ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
// defer cancel()

// var res string
// log.Printf("let's get it")

// if err := chromedp.Run(ctx,
// 	chromedp.Navigate(`https://s.cafef.vn/Lich-su-giao-dich-FPT-1.chn`),
// 	chromedp.WaitVisible(`body > h1`),
// 	// chromedp.Click(`#divSortProduct`, chromedp.NodeVisible),
// 	// chromedp.Click(`#divSortProductOptions > ul > li:nth-child(2)`, chromedp.NodeVisible),
// 	// chromedp.Click(`div > ul > li`, chromedp.NodeVisible),
// 	// chromedp.WaitVisible(`body > footer`),
// 	chromedp.Text(`h1`, &res, chromedp.NodeVisible, chromedp.ByQuery),
// ); err != nil {
// 	log.Fatal(err)
// }

// log.Printf("Go's time.After example:\n%s", res)
// }
