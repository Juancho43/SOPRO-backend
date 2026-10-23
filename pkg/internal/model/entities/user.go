package entities

import "time"

type User struct {
	UID   string
	Email string
	Name  string

	CreatedAt time.Time
	UpdatedAt time.Time
}
