package products

import (
	"fmt"
	"net/http"

	"github.com/Ricardoarsv/E-commerce_REST-API/services/auth"
	"github.com/Ricardoarsv/E-commerce_REST-API/types"
	"github.com/Ricardoarsv/E-commerce_REST-API/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductsStore
}

func NewHandler(store types.ProductsStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products/get_products", h.handleGetProduct).Methods("GET")
	router.HandleFunc("/products/create", h.handleCreateProduct).Methods("POST")
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	// todo get JWT payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// todo validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	// todo validate JWT

	err := auth.ValidateJwtTokenForCreateProducts(payload.Token)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("tokent not valid: %v", err))
		return
	}

	err = h.store.ValidateExistingProduct(&types.Products{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
	})

	if err != nil {
		if err.Error() == "product already exists" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "The Product already exists"})
		} else {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}
		return
	}

	if payload.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name cannot be empty"))
		return
	}

	err = h.store.CreateProdutc(types.Products{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
		CreatedAt:   payload.CreatedAt,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Product created"})
}
