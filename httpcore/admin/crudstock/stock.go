package crudstock

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/duongnam99/stock-analyzer/models"
	"github.com/duongnam99/stock-analyzer/repository"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseFormater struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

func Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stockAdmin models.StockAdmin
	var response ResponseFormater
	_ = json.NewDecoder(r.Body).Decode(&stockAdmin)
	if !repository.StockAdminRepository.IsStockAdminExisting(context.Background(), stockAdmin.Code) {
		err := repository.StockAdminRepository.StoreStockAdmin(stockAdmin)
		if err != nil {
			return
		}
		response = ResponseFormater{true, map[string]interface{}{"stock": stockAdmin}}
	} else {
		response = ResponseFormater{true, map[string]interface{}{"message": "Stock is already exist"}}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteByCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get params
	var params = mux.Vars(r)

	result := repository.StockAdminRepository.DeleteOneByCode(params["code"])
	response := ResponseFormater{true, map[string]interface{}{"message": result}}
	json.NewEncoder(w).Encode(response)
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
