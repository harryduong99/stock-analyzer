package crudstock

import (
	"encoding/json"
	"net/http"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func DeleteByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get params
	var params = mux.Vars(r)

	result := repository.StockAdminRepository.DeleteOneByCode(params["code"])
	json.NewEncoder(w).Encode(result)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		return
	}
	result := repository.StockAdminRepository.DeleteOneById(id)
	json.NewEncoder(w).Encode(result)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stocks := repository.StockAdminRepository.AllStockAdmin()
	json.NewEncoder(w).Encode(stocks) // encode similar to serialize process.
}
