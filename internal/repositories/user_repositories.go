package repositories

import (
	"context"
	"database/sql"

	"github.com/Nucleussss/auth-service/internal/db/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	ExistsbyEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *models.CreateNewUser) error
	FindbyEmail(ctx context.Context, email string) (*models.User, error)
	FindbyID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// ExistsbyEmail checks if a user exists by their email address.
func (ur *userRepository) ExistsbyEmail(ctx context.Context, email string) (bool, error) {
	var exist bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := ur.db.QueryRowContext(ctx, query, email).Scan(&exist)

	return exist, err
}

// Create a new user in the database.
func (ur *userRepository) Create(ctx context.Context, user *models.CreateNewUser) error {
	var query = `INSERT INTO users (name, email, password_hash) 
		VALUES ($1, $2, $3)`

	_, err := ur.db.ExecContext(ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

// FindbyEmail finds a user by their email address.
func (ur *userRepository) FindbyEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := ur.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

// FindbyID finds a user by their ID.
func (ur *userRepository) FindbyID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT *  FROM users WHERE id = $1`

	err := ur.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

// UpdatePassword updates a user's password.
func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}
