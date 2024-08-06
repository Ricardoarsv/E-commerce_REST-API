package order

import (
	"database/sql"

	"github.com/Ricardoarsv/E-commerce_REST-API/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	var id int
	query := "INSERT INTO orders (userID, total, status, address) VALUES ($1, $2, $3, $4) RETURNING id"
	err := s.db.QueryRow(query, order.UserID, order.Total, order.Status, order.Address).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) CreateOrderItem(order_item types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderID, productId, quantity, price) VALUES ($1, $2, $3, $4)", order_item.OrderID, order_item.ProductID, order_item.Quantity, order_item.Price)
	return err
}
