package cart

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
	store        types.OrderStore
	productStore types.ProductsStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductsStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods("POST")
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload

	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// todo validate the payload
	if err := utils.Validate.Struct(cart); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	ProductsIds, err := GetCartItemsIDs(cart.Items)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIds(ProductsIds)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	orderID, totalPrice, err := h.CreateOrder(ps, cart.Items, userID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
