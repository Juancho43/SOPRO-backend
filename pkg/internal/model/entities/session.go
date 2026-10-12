package entities

import "time"

// FocusmateSession representa tu bloque de tiempo de alta concentración.
type FocusmateSession struct {
	ID          string    `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"` // ej. "upcoming"
	PartnerName string    `json:"partner_name,omitempty"`
}
