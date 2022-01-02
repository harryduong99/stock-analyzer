package sourcefactory

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chromedp/chromedp"
)

type CafefSourceHandler struct {
}

func (sourceHandler CafefSourceHandler) GetData() {
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

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://s.cafef.vn/Lich-su-giao-dich-FPT-1.chn`),
		chromedp.Text(`#ctl00_ContentPlaceHolder1_ctl03_rptData2_ctl01_itemTR`, &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	infos := strings.Fields(res)

	fmt.Println(infos, len(infos))

	log.Println(strings.TrimSpace(res))
}
