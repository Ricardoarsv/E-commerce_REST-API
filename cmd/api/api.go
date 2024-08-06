package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Ricardoarsv/E-commerce_REST-API/services/cart"
	"github.com/Ricardoarsv/E-commerce_REST-API/services/order"
	"github.com/Ricardoarsv/E-commerce_REST-API/services/products"
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

	userStore := user.NewStore(api.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productsStore := products.NewStore(api.db)
	productsHandler := products.NewHandler(productsStore)
	productsHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(api.db)

	cartHandler := cart.NewHandler(orderStore, productsStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Server is running on 🚀", api.listenaddress, "(ApiServer)")

	return http.ListenAndServe(api.listenaddress, router)
}
