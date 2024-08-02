package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Ricardoarsv/E-commerce_REST-API/services/user"
	"github.com/gorilla/mux"
)

type APIserver struct {
	listenaddress string
	db            *sql.DB
}

func NewApiServer(listenaddress string, db *sql.DB) *APIserver {
	return &APIserver{
		listenaddress: listenaddress,
		db:            db,
	}
}

func (api *APIserver) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on 🚀", api.listenaddress, "(ApiServer)")

	return http.ListenAndServe(api.listenaddress, router)
}
