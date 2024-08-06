package user

import (
	"database/sql"
	"fmt"

	"github.com/Ricardoarsv/E-commerce_REST-API/types"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.DB.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no user found with email %s", email)
	}

	u, err := scanRowsIntoUser(rows)
	if err != nil {
		return nil, err
	}

	return u, nil
}
func scanRowsIntoUser(row *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAT,
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {

	rows, err := s.DB.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)

	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.DB.Exec("INSERT INTO users (firstName, lastName, email, role, password) VALUES ($1, $2, $3, $4, $5)", user.FirstName, user.LastName, user.Email, user.Role, user.Password)
	if err != nil {
		return err
	}

	return nil
}
