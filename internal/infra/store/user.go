package store

import (
	"context"
	"database/sql"
	"fmt"
	"go-web/internal/core/models"
)

func (p *pgStore) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash;
	`
	row := p.db.QueryRowContext(ctx, query, user.Id, user.Email, user.PasswordHash)
	var u models.User
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash); err != nil {
		return nil, fmt.Errorf("store.Create: %w", err)
	}
	return &u, nil
}

func (p *pgStore) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1;
	`
	row := p.db.QueryRowContext(ctx, query, email)
	var u models.User
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("store.FindByEmail: %w", err)
	}
	return &u, nil
}
