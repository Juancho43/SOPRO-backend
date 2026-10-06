package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHabitStreakRepository struct {
	db *pgxpool.Pool
}

func NewPostgresHabitStreakRepository(db *pgxpool.Pool) *PostgresHabitStreakRepository {
	return &PostgresHabitStreakRepository{db: db}
}

// Save consolida tu esfuerzo. Si la racha es nueva, la inserta; si ya existe, actualiza tu inercia.
func (r *PostgresHabitStreakRepository) Save(streak *entities.HabitStreak) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Si el ID está vacío, es el inicio de una nueva disciplina (Insert)
	if streak.ID == "" {
		query := `
			INSERT INTO habit_streaks (user_id, habit_name, frequency, current_streak, max_streak) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, created_at, updated_at`

		return r.db.QueryRow(ctx, query,
			streak.UserID,
			streak.HabitName,
			string(streak.Frequency),
			streak.CurrentStreak,
			streak.MaxStreak,
		).Scan(
			&streak.ID,
			&streak.CreatedAt,
			&streak.UpdatedAt,
		)
	}

	query := `
		UPDATE habit_streaks 
		SET current_streak = $1, max_streak = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, streak.CurrentStreak, streak.MaxStreak, streak.ID).Scan(&streak.UpdatedAt)
}

// GetByUserIDAndHabit examina tu estado actual. Es el momento de la verdad antes de sumar +1.
func (r *PostgresHabitStreakRepository) GetByUserIDAndHabit(userID string, habitName string) (*entities.HabitStreak, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, habit_name, frequency, current_streak, max_streak, created_at, updated_at 
		FROM habit_streaks 
		WHERE user_id = $1 AND habit_name = $2`

	streak := &entities.HabitStreak{}
	var freqStr string

	err := r.db.QueryRow(ctx, query, userID, habitName).Scan(
		&streak.ID,
		&streak.UserID,
		&streak.HabitName,
		&freqStr,
		&streak.CurrentStreak,
		&streak.MaxStreak,
		&streak.CreatedAt,
		&streak.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	streak.Frequency = entities.HabitFrequency(freqStr)
	return streak, nil
}

// GetAllByUserID recupera tu arsenal completo para visualizar tus conquistas en el Dashboard.
func (r *PostgresHabitStreakRepository) GetAllByUserID(userID string) ([]*entities.HabitStreak, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, habit_name, frequency, current_streak, max_streak, created_at, updated_at 
		FROM habit_streaks 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streaks []*entities.HabitStreak

	for rows.Next() {
		streak := &entities.HabitStreak{}
		var freqStr string

		if err := rows.Scan(
			&streak.ID,
			&streak.UserID,
			&streak.HabitName,
			&freqStr,
			&streak.CurrentStreak,
			&streak.MaxStreak,
			&streak.CreatedAt,
			&streak.UpdatedAt,
		); err != nil {
			return nil, err
		}

		streak.Frequency = entities.HabitFrequency(freqStr)
		streaks = append(streaks, streak)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return streaks, nil
}
