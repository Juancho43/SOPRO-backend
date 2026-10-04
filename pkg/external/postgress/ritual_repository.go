package postgress

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRitualRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRitualRepository(db *pgxpool.Pool) *PostgresRitualRepository {
	return &PostgresRitualRepository{db: db}
}

func (r *PostgresRitualRepository) ExistsForDate(userID string, date time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM daily_rituals WHERE user_id = $1 AND ritual_date = $2)`

	err := r.db.QueryRow(ctx, query, userID, date).Scan(&exists)
	return exists, err
}

func (r *PostgresRitualRepository) Save(ritual *entities.Ritual) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	goalsJSON, err := json.Marshal(ritual.Metas)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO daily_rituals (user_id, ritual_date, gratitude, goals) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at`

	err = r.db.QueryRow(ctx, query, ritual.UserID, ritual.Date, ritual.Agradecimiento, goalsJSON).Scan(
		&ritual.ID,
		&ritual.CreatedAt,
	)

	return err
}
func (r *PostgresRitualRepository) GetTodays(userID string, date time.Time) (*entities.Ritual, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, ritual_date, gratitude, goals, created_at 
		FROM daily_rituals 
		WHERE user_id = $1 AND ritual_date = $2`

	ritual := &entities.Ritual{}
	var goalsJSON []byte

	err := r.db.QueryRow(ctx, query, userID, date).Scan(
		&ritual.ID,
		&ritual.UserID,
		&ritual.Date,
		&ritual.Agradecimiento,
		&goalsJSON,
		&ritual.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if len(goalsJSON) > 0 {
		err = json.Unmarshal(goalsJSON, &ritual.Metas)
		if err != nil {
			return nil, err
		}
	}

	return ritual, nil
}
