package postgresql

import (
	"context"
	"errors"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRepository struct {
	conn *pgxpool.Pool
}

func NewLoginRepository(conn *pgxpool.Pool) *LoginRepository {
	return &LoginRepository{
		conn: conn,
	}
}

func (r LoginRepository) CreateLogin(ctx context.Context, passwordHash string) (string, error) {
	sql := `
	INSERT INTO logins (password_hash)
	VALUES ($1)
	RETURNING id
	`

	var accountId string
	err := r.conn.QueryRow(ctx, sql, passwordHash).Scan(&accountId)
	if err != nil {
		return "", err
	}

	return accountId, nil
}

func (r LoginRepository) GetPasswordHash(ctx context.Context, accountId string) (string, error) {
	sql := `
	SELECT password_hash
	FROM logins
	WHERE id = $1
	`

	var passwordHash string
	err := r.conn.QueryRow(ctx, sql, accountId).Scan(&passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", auth.ErrUserNotFound
	} else if err != nil {
		return "", err
	}

	return passwordHash, nil
}
