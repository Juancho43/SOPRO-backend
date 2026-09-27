package entities

import "time"

type User struct {
	UID   string
	Email string
	Name  string

	CurrentStreak int
	MaxStreak     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
