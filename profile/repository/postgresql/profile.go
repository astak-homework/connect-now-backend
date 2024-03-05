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
	ID        string        `db:"ID"`
	AccountID string        `db:"account_id"`
	FirstName string        `db:"first_name"`
	LastName  string        `db:"last_name"`
	BirthDate time.Time     `db:"birth_date"`
	Gender    models.Gender `db:"gender"`
	Biography string        `db:"biography"`
	City      string        `db:"city"`
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
	INSERT INTO %s (id, first_name, last_name, birth_date, gender, biography, city)
	VALUES ($1, $2, $3, $4, $5, $6, $6)
	`
	sql = fmt.Sprintf(sql, r.table)

	_, err := r.conn.Exec(ctx, sql, model.ID, model.FirstName, model.LastName, model.BirthDate, model.Gender, model.Biography, model.City)
	return err
}

func (r ProfileRepository) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	model := new(Profile)

	sql := `
	SELECT id, first_name, last_name, birth_date, gender, biography, city
	FROM %s
	WHERE id = $1
	`
	sql = fmt.Sprintf(sql, r.table)

	err := r.conn.QueryRow(ctx, sql, id).Scan(&model)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, profile.ErrProfileNotFound
	} else if err != nil {
		return nil, err
	}

	return toProfile(model), nil
}

func (r ProfileRepository) DeleteProfile(ctx context.Context, id string) error {
	sql := `
	DELETE %s
	WHERE account_id = $1
	`
	sql = fmt.Sprintf(sql, r.table)

	_, err := r.conn.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func toModel(p *models.Profile) *Profile {
	return &Profile{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		BirthDate: p.BirthDate,
		Gender:    p.Gender,
		Biography: p.Biography,
		City:      p.City,
	}
}

func toProfile(p *Profile) *models.Profile {
	return &models.Profile{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		BirthDate: p.BirthDate,
		Gender:    p.Gender,
		Biography: p.Biography,
		City:      p.City,
	}
}
