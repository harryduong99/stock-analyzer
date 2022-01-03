package crawler

import (
	"github.com/duongnam99/stock-analyzer/crawler/sourcefactory"
)

func Crawl(source string) {
	targets := getTargets()
	sourcefactory.GetSourceHandlerFactory(source).GetData(targets)
}

func getTargets() []string {
	return []string{"FPT", "FLC", "VPB", "TCB", "VCB", "HPG"}
}
