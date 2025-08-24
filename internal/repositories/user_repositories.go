package repositories

import (
	"database/sql"

	"github.com/Nucleussss/auth-service/internal/db/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsbyEmail(email string) (bool, error)
	Create(name, email, passwordHash string) error
	FindbyEmail(email string) (*models.User, error)
	FindbyID(id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// ExistsbyEmail checks if a user exists by their email address.
func (ur *userRepository) ExistsbyEmail(email string) (bool, error) {
	var exist bool
	err := ur.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email,
	).Scan(&exist)

	return exist, err
}

// Create creates a new user in the database.
func (ur *userRepository) Create(name, email, passwordHash string) error {
	_, err := ur.db.Exec(
		`INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3)`,
		name,
		email,
		passwordHash,
	)

	return err
}

// FindbyEmail finds a user by their email address.
func (ur *userRepository) FindbyEmail(email string) (*models.User, error) {
	var user models.User

	err := ur.db.QueryRow(
		`SELECT id, name, email, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email)

	return &user, err
}

// FindbyID finds a user by their ID.
func (ur *userRepository) FindbyID(userID uuid.UUID) (*models.User, error) {
	var user models.User

	err := ur.db.QueryRow(
		`SELECT id, name, email, password_hash FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Name, &user.Email)

	return &user, err
}
