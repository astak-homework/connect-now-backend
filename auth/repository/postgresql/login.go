package postgresql

import (
	"context"

	"github.com/astak-homework/connect-now-backend/auth"
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

func (r LoginRepository) AuthenticateLogin(ctx context.Context, accountId, passwordHash string) error {
	sql := `
	SELECT id, user_name, password_hash
	FROM logins
	WHERE user_name = $1 and password_hash = $2
	`

	rows, err := r.conn.Query(ctx, sql, accountId, passwordHash)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return auth.ErrUserNotFound
	}

	return nil
}
