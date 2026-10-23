package entities

import "time"

// Habit representa un hábito atómico en el arsenal del usuario.
type Habit struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Frequency string    `json:"frequency"` // "daily" | "weekly"
	Archived  bool      `json:"archived"`
	DomainID  *string   `json:"domain_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HabitLog registra el cumplimiento individual de un hábito en una fecha calendario.
type HabitLog struct {
	ID          string    `json:"id"`
	HabitID     string    `json:"habit_id"`
	UserID      string    `json:"user_id"`
	LogDate     string    `json:"log_date"` // YYYY-MM-DD
	CompletedAt time.Time `json:"completed_at"`
}

// HabitWithStreak combina la definición del hábito con su racha actual y estado del día.
type HabitWithStreak struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Frequency      string    `json:"frequency"`
	Archived       bool      `json:"archived"`
	DomainID       *string   `json:"domain_id,omitempty"`
	DomainName     *string   `json:"domain_name,omitempty"`
	DomainColor    *string   `json:"domain_color,omitempty"`
	CurrentStreak  int       `json:"current_streak"`
	MaxStreak      int       `json:"max_streak"`
	CompletedToday bool      `json:"completed_today"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// HabitCheckInResult representa la respuesta tras completar un hábito.
type HabitCheckInResult struct {
	HabitID       string `json:"habit_id"`
	LogDate       string `json:"log_date"`
	CurrentStreak int    `json:"current_streak"`
	MaxStreak     int    `json:"max_streak"`
	Message       string `json:"message"`
}
