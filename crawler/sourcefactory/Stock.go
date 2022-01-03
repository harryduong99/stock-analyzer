package sourcefactory

type StockInfo struct {
	Date             string
	AdjustedPrice    float32
	ClosedPrice      float32
	Change           string
	StockOrderAmount float32
	StockOrderValue  float32
	StockDealAmount  float32
	StockDealValue   float32
	OpenPrice        float32
	HighestPrice     float32
	LowestPrice      float32
}
