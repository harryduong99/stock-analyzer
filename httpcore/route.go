package httpcore

import (
	"log"
	"net/http"

	"github.com/duongnam99/stock-analyzer/httpcore/admin"
	"github.com/duongnam99/stock-analyzer/httpcore/admin/crudstock"
	"github.com/duongnam99/stock-analyzer/httpcore/middleware"
	"github.com/gorilla/mux"
)

func InitRoutes() {
	r := mux.NewRouter()

	apiRoute := r.PathPrefix("/api/").Subrouter()
	apiRoute.HandleFunc("/v1/stocks", crudstock.GetAll).Methods("GET")
	apiRoute.HandleFunc("/v1/stocks", crudstock.Store).Methods("POST")
	apiRoute.HandleFunc("/v1/stocks/{code}", crudstock.DeleteByCode).Methods("DELETE")
	apiRoute.Use(middleware.AdminMiddleware())
	r.HandleFunc("/stock/report", admin.Report).Methods("GET")

	log.Fatal(http.ListenAndServe(":6868", r))
}
