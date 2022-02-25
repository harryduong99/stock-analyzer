# Instruction

## Run

- To crawl all information for last 3 days, then analyze and send the email, run:

`go run main.go analyze 3`

- To crawl stock data of the last 10 days, run:

`go run main.go crawl cafef 10`

### Supporting sources

- `cafef`

- `vietstock`

- In the future, we will have `vndirect`

- Crawling vietstock need to use headless browser so we will crawl cafef first, then vietstock will be crawled to fulfill the data that cafef is missing (or not)