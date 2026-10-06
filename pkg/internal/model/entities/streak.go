package entities

import "time"

type HabitStreak struct {
	ID            string
	UserID        string
	HabitName     string
	Frequency     HabitFrequency
	CurrentStreak int
	MaxStreak     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type HabitFrequency string

const (
	// Daily representa la disciplina inquebrantable de cada día.
	Daily HabitFrequency = "diario"

	// Weekly enfoca tu energía en victorias consolidadas a lo largo de la semana.
	Weekly HabitFrequency = "semanal"
)

// IsValid actúa como un guardián: rechaza cualquier entrada que no cumpla tus estándares.
func (f HabitFrequency) IsValid() bool {
	switch f {
	case Daily, Weekly:
		return true
	default:
		return false
	}
}
