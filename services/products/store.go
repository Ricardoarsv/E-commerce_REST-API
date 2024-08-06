package products

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Ricardoarsv/E-commerce_REST-API/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateProdutc(product types.Products) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, quantity, price, create_at) VALUES ($1, $2, $3, $4, $5, $6)", product.Name, product.Description, product.Image, product.Quantity, product.Price, product.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProducts() ([]*types.Products, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Products, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) ValidateExistingProduct(product *types.Products) error {
	rows, err := s.db.Query("SELECT * FROM products WHERE name = $1 AND description = $2 AND price = $3", product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return fmt.Errorf("product already exists")
	}

	return nil
}

func (s *Store) GetProductsByIds(productIDs []int) ([]types.Products, error) {
	placeholders := make([]string, len(productIDs))
	for i := range productIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", strings.Join(placeholders, ","))

	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.Products{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Store) UpdateProductStock(product types.Products) error {
	_, err := s.db.Exec("UPDATE products SET quantity = $1 WHERE id = $2", product.Quantity, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Products, error) {
	product := new(types.Products)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Quantity,
		&product.Price,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
