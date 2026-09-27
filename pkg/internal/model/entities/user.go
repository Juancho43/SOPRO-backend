package entities

import "time"

type User struct {
	UID   string // Must be uppercase 'ID', not 'id'
	Email string
	Name  string

	CurrentStreak int       // Must be uppercase 'C'
	MaxStreak     int       // Must be uppercase 'M'
	CreatedAt     time.Time // Must be uppercase 'C'
	UpdatedAt     time.Time // Must be uppercase 'U'
}
