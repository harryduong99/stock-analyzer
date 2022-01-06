package crudstock

import (
	"log"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
)

func Add(stock string) {
	stockAdmin := models.StockAdmin{}
	stockAdmin.Code = stock
	err := repository.StockAdminRepository.StoreStockAdmin(stockAdmin)
	if err != nil {
		log.Println(stock)
	}
}

func getAll() []models.StockAdmin {
	stocks := repository.StockAdminRepository.AllStockAdmin()
	return stocks
}
