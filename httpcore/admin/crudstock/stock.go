package crudstock

import (
	"encoding/json"
	"net/http"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
)

func Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stockAdmin models.StockAdmin

	_ = json.NewDecoder(r.Body).Decode(&stockAdmin)
	err := repository.StockAdminRepository.StoreStockAdmin(stockAdmin)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(stockAdmin)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stocks := repository.StockAdminRepository.AllStockAdmin()
	json.NewEncoder(w).Encode(stocks) // encode similar to serialize process.
}
