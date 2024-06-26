package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/astak-homework/connect-now-backend/models"
	"github.com/astak-homework/connect-now-backend/profile"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID        string        `db:"id"`
	FirstName string        `db:"first_name"`
	LastName  string        `db:"last_name"`
	BirthDate time.Time     `db:"birth_date"`
	Gender    models.Gender `db:"gender"`
	Biography string        `db:"biography"`
	City      string        `db:"city"`
}

type ProfileRepository struct {
	conn *pgxpool.Pool
}

func NewProfileRepository(conn *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		conn: conn,
	}
}

func (r ProfileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	model := toModel(profile)

	sql := `
	INSERT INTO profiles (id, first_name, last_name, birth_date, gender, biography, city)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.conn.Exec(ctx, sql, model.ID, model.FirstName, model.LastName, model.BirthDate, model.Gender, model.Biography, model.City)
	return err
}

func (r ProfileRepository) GetProfile(ctx context.Context, id string) (*models.Profile, error) {
	model := new(Profile)

	sql := `
	SELECT id, first_name, last_name, birth_date, gender, biography, city
	FROM profiles
	WHERE id = $1
	`

	err := pgxscan.Get(ctx, r.conn, model, sql, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, profile.ErrProfileNotFound
	} else if err != nil {
		return nil, err
	}

	return toProfile(model), nil
}

func (r ProfileRepository) DeleteProfile(ctx context.Context, id string) error {
	sql := `
	DELETE profiles
	WHERE account_id = $1
	`

	_, err := r.conn.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (r ProfileRepository) SearchProfile(ctx context.Context, firstName, lastName string) ([]*models.Profile, error) {
	sql := `
	SELECT id, first_name, last_name, birth_date, gender, biography, city
	FROM profiles
	WHERE first_name like $1 and last_name like $2
	ORDER BY id
	`

	var profiles []*Profile
	err := pgxscan.Select(ctx, r.conn, &profiles, sql, firstName+"%", lastName+"%")

	if err != nil {
		return nil, err
	}

	return toProfiles(profiles), nil
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

func toProfiles(ps []*Profile) []*models.Profile {
	profiles := []*models.Profile{}
	for _, p := range ps {
		profiles = append(profiles, toProfile(p))
	}
	return profiles
}
