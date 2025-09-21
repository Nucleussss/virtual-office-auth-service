package repositories

import (
	"context"
	"database/sql"

	"github.com/Nucleussss/auth-service/internal/db/models"
)

type PasswordResetRepository interface {
	Create(ctx context.Context, reset *models.PasswordReset) error
	FindValidToken(ctx context.Context, token string) (*models.PasswordReset, error)
	Delete(ctx context.Context, token string) error
	DeleteExpiresTokens(ctx context.Context) error
}

type passwordResetRepository struct {
	db *sql.DB
}

func NewPasswordRepository(db *sql.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// Create a new password reset record in the database.
func (r *passwordResetRepository) Create(ctx context.Context, reset *models.PasswordReset) error {
	query := ` 
		INSERT INTO password_resets (token, user_id, expired_at) 
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query,
		reset.Token,
		reset.UserID,
		reset.ExpiresAt,
	)
	return err
}

// Find a valid password reset token in the database.
func (r *passwordResetRepository) FindValidToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	query := `
		SELECT * FROM password_resets 
		WHERE token = $1 && expired_at >= NOW()
	`

	row := r.db.QueryRowContext(ctx, query, token)

	var reset *models.PasswordReset
	err := row.Scan(
		&reset.Token,
		&reset.UserID,
		&reset.ExpiresAt,
		&reset.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return reset, nil
}

func (r *passwordResetRepository) Delete(ctx context.Context, token string) error {
	query := `
		DELETE * FROM password_resets
		WHERE token = $1
	`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *passwordResetRepository) DeleteExpiresTokens(ctx context.Context) error {
	query := `
		DELETE * FROM password_resets
		WHERE expired_at <= NOW()
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
