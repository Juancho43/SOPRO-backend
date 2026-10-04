package entities

import "time"

type Ritual struct {
	ID             string
	UserID         string
	Date           time.Time
	Agradecimiento string
	Metas          []string // El éxito exige enfoque: exactamente 10 metas
	CreatedAt      time.Time
}
