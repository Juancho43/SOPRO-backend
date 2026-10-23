package entities

import "time"

// Domain representa una dimensión de excelencia o categoría de desarrollo personal.
type Domain struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

// DomainEffortPoint representa el punto de datos para el radar de esfuerzo semanal.
type DomainEffortPoint struct {
	DomainID       string  `json:"domain_id"`
	DomainName     string  `json:"domain_name"`
	Color          string  `json:"color"`
	CompletedCount int     `json:"completed_count"`
	EffortScore    float64 `json:"effort_score"` // Normalizado de 0 a 100
}

// DomainWithHabitCount incluye el dominio junto con la cantidad de hábitos asociados.
type DomainWithHabitCount struct {
	Domain
	HabitCount int `json:"habit_count"`
}
