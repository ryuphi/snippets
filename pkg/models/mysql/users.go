package mysql

import (
	"database/sql"
	"learn-web/snippets/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert method add a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate method to verify whether a user exists with the provided
// email address ans password. This will return the relevant user id if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get method to fetch details for a specific user based
func (m UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
