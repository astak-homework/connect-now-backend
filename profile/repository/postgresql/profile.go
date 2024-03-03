package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID        string    `db:"ID"`
	AccountID string    `db:"account_id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	BirthDate time.Time `db:"birth_date"`
	Gender    string    `db:"gender"`
	Biography string    `db:"biography"`
	City      string    `db:"city"`
}

type ProfileRepository struct {
	conn  *pgxpool.Pool
	table string
}

func NewProfileRepository(conn *pgxpool.Pool, tableName string) *ProfileRepository {
	return &ProfileRepository{
		conn:  conn,
		table: tableName,
	}
}

func (r ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	model := toModel(profile)

	sql := `
	INSERT INTO %s (account_id, first_name, last_name, birth_date, gender, biography, city)
	VALUES ($1, $2, $3, $4, $5, $6, $6)
	RETURNING id
	`
	sql = fmt.Sprintf(sql, r.table)

	err := r.conn.QueryRow(ctx, sql, profile.Account.ID, profile.FirstName, profile.LastName, profile.BirthDate, profile.Gender, profile.Biography, profile.City).Scan(&model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r ProfileRepository) GetProfile(ctx context.Context, account *models.Login) (*models.Profile, error) {
	model := new(Profile)

	sql := `
	SELECT id, account_id, first_name, last_name, birth_date, gender, biography, city
	FROM %s
	WHERE account_id = $1
	`
	sql = fmt.Sprintf(sql, r.table)

	err := r.conn.QueryRow(ctx, sql, account.ID).Scan(&model)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, profile.ErrProfileNotFound
	} else if err != nil {
		return nil, err
	}

	return toProfile(model, account), nil
}

func (r ProfileRepository) DeleteProfile(ctx context.Context, account *models.Login) error {
	sql := `
	DELETE %s
	WHERE account_id = $1
	`
	sql = fmt.Sprintf(sql, r.table)

	_, err := r.conn.Exec(ctx, sql, account.ID)
	if err != nil {
		return err
	}

	return nil
}

func toModel(p *models.Profile) *Profile {
	return &Profile{
		ID:        p.ID,
		AccountID: p.Account.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		BirthDate: p.BirthDate,
		Gender:    string(p.Gender),
		Biography: p.Biography,
		City:      p.City,
	}
}

func toProfile(p *Profile, account *models.Login) *models.Profile {
	return &models.Profile{
		ID:        p.ID,
		Account:   account,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		BirthDate: p.BirthDate,
		Gender:    models.Gender(p.Gender),
		Biography: p.Biography,
		City:      p.City,
	}
}
