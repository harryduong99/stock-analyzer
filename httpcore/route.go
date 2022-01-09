package httpcore

import (
	"log"
	"net/http"

	"github.com/duongnam99/stock-analyzer/httpcore/admin"
	"github.com/duongnam99/stock-analyzer/httpcore/admin/crudstock"
	"github.com/gorilla/mux"
)

func InitRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/stocks", crudstock.GetAll).Methods("GET")
	// r.HandleFunc("/api/stocks/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/v1/stocks", crudstock.Store).Methods("POST")
	// r.HandleFunc("/api/stocks/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/api/v1/stocks/{id}", crudstock.Delete).Methods("DELETE")
	r.HandleFunc("/api/v1/stocks/{code}", crudstock.DeleteByCode).Methods("DELETE")
	r.HandleFunc("/api/v1/report", admin.Report).Methods("GET")

	log.Fatal(http.ListenAndServe(":6868", r))
}
