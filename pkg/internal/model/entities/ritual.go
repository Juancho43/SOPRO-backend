package entities

import "time"

type Ritual struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Date           time.Time `json:"date"`
	Agradecimiento string    `json:"grateful_for"`
	Metas          []string  `json:"goals"`
	CreatedAt      time.Time `json:"created_at"`
}
