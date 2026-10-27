package entities

import "time"

// ChallengeStatus representa los estados posibles de un reto estratégico.
type ChallengeStatus string

const (
	ChallengeStatusActive    ChallengeStatus = "ACTIVE"
	ChallengeStatusCompleted ChallengeStatus = "COMPLETED"
	ChallengeStatusFailed    ChallengeStatus = "FAILED"
	ChallengeStatusArchived  ChallengeStatus = "ARCHIVED"
	ChallengeStatusCancelled ChallengeStatus = "CANCELLED"
)

// Challenge representa un reto estratégico a plazo fijo (ej. 30, 60, 90 días).
type Challenge struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	StartDate   string          `json:"start_date"` // YYYY-MM-DD
	EndDate     string          `json:"end_date"`   // YYYY-MM-DD
	Reward      string          `json:"reward"`
	Status      ChallengeStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ChallengeDetail incluye el reto y el cálculo de avance basado en sus hitos.
type ChallengeDetail struct {
	Challenge
	Milestones       []*Milestone `json:"milestones"`
	TotalMilestones  int          `json:"total_milestones"`
	DoneMilestones   int          `json:"done_milestones"`
	ProgressPercent  float64      `json:"progress_percent"`
	DaysRemaining    int          `json:"days_remaining"`
	TotalDays        int          `json:"total_days"`
	DaysElapsed      int          `json:"days_elapsed"`
}

// Milestone representa una deadline intermedia o hito dentro de un reto activo.
type Milestone struct {
	ID          string     `json:"id"`
	ChallengeID string     `json:"challenge_id"`
	Title       string     `json:"title"`
	DueDate     string     `json:"due_date"` // YYYY-MM-DD
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	OrderIndex  int        `json:"order_index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UpcomingDeadline alerta visual y de vencimiento para el Dashboard.
type UpcomingDeadline struct {
	MilestoneID    string `json:"milestone_id"`
	ChallengeID    string `json:"challenge_id"`
	ChallengeTitle string `json:"challenge_title"`
	Title          string `json:"title"`
	DueDate        string `json:"due_date"`
	DaysRemaining  int    `json:"days_remaining"`
	UrgencyLevel   string `json:"urgency_level"` // "overdue", "critical" (<=1d), "warning" (<=3d), "normal"
}
