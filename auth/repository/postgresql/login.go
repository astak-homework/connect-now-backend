package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/astak-homework/connect-now-backend/auth"
	"github.com/astak-homework/connect-now-backend/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Login struct {
	ID       string `db:"id"`
	UserName string `db:"user_name"`
	Password string `db:"password_hash"`
}

type LoginRepository struct {
	conn  *pgxpool.Pool
	table string
}

func NewLoginRepository(conn *pgxpool.Pool, tableName string) *LoginRepository {
	return &LoginRepository{
		conn:  conn,
		table: tableName,
	}
}

func (r LoginRepository) CreateLogin(ctx context.Context, login *models.Login) error {
	model := toModel(login)

	sql := `
	INSERT INTO %s (user_name, password_hash)
	VALUES ($1, $2)
	`
	sql = fmt.Sprintf(sql, r.table)

	_, err := r.conn.Exec(ctx, sql, model.UserName, model.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r LoginRepository) GetLogin(ctx context.Context, username, password string) (*models.Login, error) {
	model := new(Login)

	sql := `
	SELECT id, user_name, password_hash
	FROM %s
	WHERE user_name = $1 and password_hash = $2
	`
	sql = fmt.Sprintf(sql, r.table)

	rows, err := r.conn.Query(ctx, sql, username, password)
	if err != nil {
		return nil, err
	}

	err = pgxscan.ScanOne(model, rows)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, auth.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return toProfile(model), nil
}

func toModel(l *models.Login) *Login {
	return &Login{
		UserName: l.UserName,
		Password: l.Password,
	}
}

func toProfile(l *Login) *models.Login {
	return &models.Login{
		ID:       l.ID,
		UserName: l.UserName,
		Password: l.Password,
	}
}
