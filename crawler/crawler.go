package crawler

import (
	"github.com/duongnam99/stock-analyzer/crawler/sourcefactory"
)

func Crawl(source string, totalDays int) {
	targets := getTargets()
	sourcefactory.GetSourceHandlerFactory(source).GetData(targets, totalDays, "")
}

func getTargets() []string {
	// return []string{"FPT", "FLC", "VPB", "TCB", "VCB", "HPG"}
	return []string{"FPT"}
}
