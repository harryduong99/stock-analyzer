package main

import (
	"os"

	"github.com/duongnam99/stock-analyzer/crawler"
)

func main() {
	crawler.Crawl(getSource())
}

func getSource() string {
	arg := os.Args[1:2]
	source := arg[0]
	return source
}
