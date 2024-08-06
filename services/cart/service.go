package cart

import (
	"fmt"

	"github.com/Ricardoarsv/E-commerce_REST-API/types"
)

func GetCartItemsIDs(items []types.CartItem) ([]int, error) {
	productsIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productsIds[i] = item.ProductID
	}

	return productsIds, nil
}

func (h *Handler) CreateOrder(ps []types.Products, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Products)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	totalprice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity = product.Quantity - item.Quantity
		productMap[item.ProductID] = product
		h.productStore.UpdateProductStock(product)
	}

	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalprice,
		Status:  "pending",
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalprice, nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Products) float64 {
	var totalPrice float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		totalPrice += product.Price * float64(item.Quantity)
	}

	return totalPrice
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Products) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product := products[item.ProductID]

		if product.Quantity <= 0 {
			return fmt.Errorf("product %s is not avaible in the store", product.Name)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not avaible in the quantity requested", product.Name)
		}
	}

	return nil
}
