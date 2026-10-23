package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresHabitRepository struct {
	db *pgxpool.Pool
}

func NewPostgresHabitRepository(db *pgxpool.Pool) *PostgresHabitRepository {
	return &PostgresHabitRepository{db: db}
}

func (r *PostgresHabitRepository) Create(ctx context.Context, habit *entities.Habit) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO habits (user_id, name, frequency, domain_id) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, archived, created_at, updated_at`

	err = tx.QueryRow(ctx, query, habit.UserID, habit.Name, habit.Frequency, habit.DomainID).Scan(
		&habit.ID,
		&habit.Archived,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Inicializamos el registro en habit_streaks
	streakQuery := `
		INSERT INTO habit_streaks (user_id, habit_name, habit_id, frequency, current_streak, max_streak, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 0, 0, $5, $5)`

	_, err = tx.Exec(ctx, streakQuery, habit.UserID, habit.Name, habit.ID, habit.Frequency, habit.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresHabitRepository) GetByID(ctx context.Context, id string) (*entities.Habit, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, name, frequency, archived, domain_id, created_at, updated_at 
		FROM habits 
		WHERE id = $1`

	habit := &entities.Habit{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&habit.ID,
		&habit.UserID,
		&habit.Name,
		&habit.Frequency,
		&habit.Archived,
		&habit.DomainID,
		&habit.CreatedAt,
		&habit.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return habit, nil
}

func (r *PostgresHabitRepository) ListByUserID(ctx context.Context, userID string, todayDate string) ([]*entities.HabitWithStreak, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			h.id, h.user_id, h.name, h.frequency, h.archived, h.domain_id,
			d.name AS domain_name, d.color AS domain_color,
			COALESCE(hs.current_streak, 0) AS current_streak,
			COALESCE(hs.max_streak, 0) AS max_streak,
			EXISTS (
				SELECT 1 FROM habit_logs hl 
				WHERE hl.habit_id = h.id AND hl.log_date = $2
			) AS completed_today,
			h.created_at, h.updated_at
		FROM habits h
		LEFT JOIN domains d ON d.id = h.domain_id
		LEFT JOIN habit_streaks hs ON hs.habit_id = h.id
		WHERE h.user_id = $1 AND h.archived = FALSE
		ORDER BY h.created_at ASC`

	rows, err := r.db.Query(ctx, query, userID, todayDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []*entities.HabitWithStreak
	for rows.Next() {
		item := &entities.HabitWithStreak{}
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Frequency,
			&item.Archived,
			&item.DomainID,
			&item.DomainName,
			&item.DomainColor,
			&item.CurrentStreak,
			&item.MaxStreak,
			&item.CompletedToday,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		habits = append(habits, item)
	}

	return habits, rows.Err()
}

func (r *PostgresHabitRepository) Update(ctx context.Context, habit *entities.Habit) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE habits 
		SET name = $1, frequency = $2, domain_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5`

	cmdTag, err := r.db.Exec(ctx, query, habit.Name, habit.Frequency, habit.DomainID, habit.ID, habit.UserID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("hábito no encontrado o no pertenece al usuario")
	}

	return nil
}

func (r *PostgresHabitRepository) Archive(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE habits SET archived = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("hábito no encontrado o no pertenece al usuario")
	}

	return nil
}

func (r *PostgresHabitRepository) LogCompletion(ctx context.Context, log *entities.HabitLog) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO habit_logs (habit_id, user_id, log_date, completed_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return r.db.QueryRow(ctx, query, log.HabitID, log.UserID, log.LogDate, log.CompletedAt).Scan(&log.ID)
}

func (r *PostgresHabitRepository) HasLogForDate(ctx context.Context, habitID string, date string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = $1 AND log_date = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, habitID, date).Scan(&exists)
	return exists, err
}

func (r *PostgresHabitRepository) GetLogs(ctx context.Context, habitID string, fromDate string, toDate string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT TO_CHAR(log_date, 'YYYY-MM-DD') 
		FROM habit_logs 
		WHERE habit_id = $1 AND log_date >= $2 AND log_date <= $3 
		ORDER BY log_date ASC`

	rows, err := r.db.Query(ctx, query, habitID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}

	return dates, rows.Err()
}

func (r *PostgresHabitRepository) GetStreakByHabitID(ctx context.Context, habitID string) (int, int, *time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT current_streak, max_streak, updated_at 
		FROM habit_streaks 
		WHERE habit_id = $1`

	var currentStreak, maxStreak int
	var updatedAt time.Time

	err := r.db.QueryRow(ctx, query, habitID).Scan(&currentStreak, &maxStreak, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, nil, nil
		}
		return 0, 0, nil, err
	}

	return currentStreak, maxStreak, &updatedAt, nil
}

func (r *PostgresHabitRepository) SaveHabitStreak(
	ctx context.Context,
	habitID string,
	userID string,
	habitName string,
	frequency string,
	currentStreak int,
	maxStreak int,
	now time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Primero verificamos si ya existe por habit_id
	checkQuery := `SELECT id FROM habit_streaks WHERE habit_id = $1`
	var existingID string
	err := r.db.QueryRow(ctx, checkQuery, habitID).Scan(&existingID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if existingID != "" {
		updateQuery := `
			UPDATE habit_streaks 
			SET current_streak = $1, max_streak = $2, updated_at = $3 
			WHERE id = $4`
		_, err = r.db.Exec(ctx, updateQuery, currentStreak, maxStreak, now, existingID)
		return err
	}

	insertQuery := `
		INSERT INTO habit_streaks (habit_id, user_id, habit_name, frequency, current_streak, max_streak, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $7)`
	_, err = r.db.Exec(ctx, insertQuery, habitID, userID, habitName, frequency, currentStreak, maxStreak, now)
	return err
}
