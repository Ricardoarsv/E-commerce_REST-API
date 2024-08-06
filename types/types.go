package types

import "time"

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Role      int    `json:"role"  validate:"required"`
	Password  string `json:"password" validate:"required,min=8,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

type CreateProductPayload struct {
	Token       string    `json:"token" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type ProductsStore interface {
	GetProducts() ([]*Products, error)
	GetProductsByIds(ps []int) ([]Products, error)
	CreateProdutc(Products) error
	ValidateExistingProduct(*Products) error
	UpdateProductStock(Products) error
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type Order struct {
	ID      int     `json:"id"`
	UserID  int     `json:"userID"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
	Address string  `json:"address"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreateAt  time.Time `json:"createAt"`
}

type Products struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type User struct {
	ID        int       `json:"id"`
	Role      int       `json:"role"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAT time.Time `json:"createdAt"`
}

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}
