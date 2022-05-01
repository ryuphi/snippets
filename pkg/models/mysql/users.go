package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"learn-web/snippets/pkg/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

// Insert method add a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	// create bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	stmt := `insert into users (name, email, hashed_password, created)
	values(?,?,?, utc_timestamp())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	if err != nil {
		// if this return an error, we use the error.As() function to check whether the error
		// has the type *mysql.MySQLError. If it does, the error will be assigned to the mySQLError
		// variable.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
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
