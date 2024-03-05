package models

import "time"

type Gender string

const (
	GenderMale  = Gender("male")
	GnderFemale = Gender("female")
)

type Profile struct {
	ID        string
	FirstName string
	LastName  string
	BirthDate time.Time
	Gender    Gender
	Biography string
	City      string
}
