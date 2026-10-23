package entities

import "time"

type TrophyTargetType string

const (
	TargetMasterStreak TrophyTargetType = "master_streak"
	TargetHabit        TrophyTargetType = "habit"
	TargetSynergy      TrophyTargetType = "synergy"
)

// Trophy define la recompensa u objetivo gamificado.
type Trophy struct {
	ID                 string           `json:"id"`
	UserID             string           `json:"user_id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	Icon               string           `json:"icon"`
	Color              string           `json:"color"`
	TargetType         TrophyTargetType `json:"target_type"`
	HabitID            *string          `json:"habit_id,omitempty"`
	RequiredStreakDays int              `json:"required_streak_days"`
	CreatedAt          time.Time        `json:"created_at"`
}

// TrophyCondition representa una condición AND para un trofeo de sinergia multihábito.
type TrophyCondition struct {
	ID                 string    `json:"id"`
	TrophyID           string    `json:"trophy_id"`
	HabitID            string    `json:"habit_id"`
	RequiredStreakDays int       `json:"required_streak_days"`
	CreatedAt          time.Time `json:"created_at"`
}

// UserTrophy representa la vinculación de un trofeo desbloqueado por un usuario.
type UserTrophy struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TrophyID   string    `json:"trophy_id"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

// ShowcaseTrophy representa el trofeo en la vitrina del usuario con su estado actual de desbloqueo.
type ShowcaseTrophy struct {
	ID                 string            `json:"id"`
	UserID             string            `json:"user_id"`
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Icon               string            `json:"icon"`
	Color              string            `json:"color"`
	TargetType         TrophyTargetType  `json:"target_type"`
	HabitID            *string           `json:"habit_id,omitempty"`
	RequiredStreakDays int               `json:"required_streak_days"`
	Unlocked           bool              `json:"unlocked"`
	UnlockedAt         *time.Time        `json:"unlocked_at,omitempty"`
	Conditions         []TrophyCondition `json:"conditions,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

// SynergyHabitProgress detalla el avance de un hábito requerido en un trofeo de sinergia.
type SynergyHabitProgress struct {
	HabitID            string `json:"habit_id"`
	HabitName          string `json:"habit_name"`
	CurrentStreak      int    `json:"current_streak"`
	RequiredStreakDays int    `json:"required_streak_days"`
	IsMet              bool   `json:"is_met"`
}

// SynergyStatusResponse consolida el estado de progreso de un trofeo de sinergia.
type SynergyStatusResponse struct {
	TrophyID   string                 `json:"trophy_id"`
	TrophyName string                 `json:"trophy_name"`
	Unlocked   bool                   `json:"unlocked"`
	Progress   []SynergyHabitProgress `json:"progress"`
}

// HabitCompletedEvent es el payload emitido en el bus cuando se checkea un hábito.
type HabitCompletedEvent struct {
	UserID     string `json:"user_id"`
	HabitID    string `json:"habit_id"`
	HabitName  string `json:"habit_name"`
	StreakDays int    `json:"streak_days"`
	Date       string `json:"date"`
}

// TrophyUnlockedEvent es el payload emitido cuando se desbloquea un trofeo.
type TrophyUnlockedEvent struct {
	UserID     string    `json:"user_id"`
	TrophyID   string    `json:"trophy_id"`
	TrophyName string    `json:"trophy_name"`
	Icon       string    `json:"icon"`
	Color      string    `json:"color"`
	UnlockedAt time.Time `json:"unlocked_at"`
	IsSynergy  bool      `json:"is_synergy"`
}
