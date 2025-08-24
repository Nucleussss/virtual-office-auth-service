package repositories

import "database/sql"

type PasswordResetRepository interface {
}

type passwordResetRepository struct {
	db *sql.DB
}

func NewPasswordRepository(db *sql.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}
