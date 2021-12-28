package crawler

import (
	"github.com/duongnam99/stock-analyzer/crawler/sourcefactory"
)

func Crawl(source string) {
	sourcefactory.GetSourceHandlerFactory(source).GetData()
}
