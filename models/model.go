package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockInfo struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // tag golang
	Date             string             `json:"date" bson:"date"`
	AdjustedPrice    float64            `json:"adjusted_price" bson:"adjusted_price"`
	ClosedPrice      float64            `json:"closed_price" bson:"closed_price"`
	Change           string             `json:"change" bson:"change"`
	StockOrderAmount int                `json:"stock_order_amount" bson:"stock_order_amount"`
	StockOrderValue  int                `json:"stock_order_value" bson:"stock_order_value"`
	StockDealAmount  int                `json:"stock_deal_amount" bson:"stock_deal_amount"`
	StockDealValue   int                `json:"stock_deal_value" bson:"stock_deal_value"`
	OpenPrice        float64            `json:"open_price" bson:"open_price"`
	HighestPrice     float64            `json:"higest_price" bson:"higest_price"`
	LowestPrice      float64            `json:"lowest_price" bson:"lowest_price"`
}
