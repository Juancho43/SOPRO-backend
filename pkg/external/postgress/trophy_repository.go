package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTrophyRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTrophyRepository(db *pgxpool.Pool) *PostgresTrophyRepository {
	return &PostgresTrophyRepository{db: db}
}

func (r *PostgresTrophyRepository) Create(ctx context.Context, trophy *entities.Trophy, conditions []entities.TrophyCondition) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO trophies (user_id, name, description, icon, color, target_type, habit_id, required_streak_days)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err = tx.QueryRow(ctx, query,
		trophy.UserID,
		trophy.Name,
		trophy.Description,
		trophy.Icon,
		trophy.Color,
		string(trophy.TargetType),
		trophy.HabitID,
		trophy.RequiredStreakDays,
	).Scan(&trophy.ID, &trophy.CreatedAt)
	if err != nil {
		return err
	}

	if trophy.TargetType == entities.TargetSynergy && len(conditions) > 0 {
		condQuery := `
			INSERT INTO trophy_conditions (trophy_id, habit_id, required_streak_days)
			VALUES ($1, $2, $3)`

		for _, cond := range conditions {
			_, err = tx.Exec(ctx, condQuery, trophy.ID, cond.HabitID, cond.RequiredStreakDays)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresTrophyRepository) GetByID(ctx context.Context, id string) (*entities.Trophy, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, name, description, icon, color, target_type, habit_id, required_streak_days, created_at
		FROM trophies 
		WHERE id = $1`

	trophy := &entities.Trophy{}
	var targetTypeStr string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&trophy.ID,
		&trophy.UserID,
		&trophy.Name,
		&trophy.Description,
		&trophy.Icon,
		&trophy.Color,
		&targetTypeStr,
		&trophy.HabitID,
		&trophy.RequiredStreakDays,
		&trophy.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	trophy.TargetType = entities.TrophyTargetType(targetTypeStr)
	return trophy, nil
}

func (r *PostgresTrophyRepository) ListByUserID(ctx context.Context, userID string) ([]*entities.Trophy, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, name, description, icon, color, target_type, habit_id, required_streak_days, created_at
		FROM trophies 
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trophies []*entities.Trophy
	for rows.Next() {
		t := &entities.Trophy{}
		var targetTypeStr string
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Name,
			&t.Description,
			&t.Icon,
			&t.Color,
			&targetTypeStr,
			&t.HabitID,
			&t.RequiredStreakDays,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		t.TargetType = entities.TrophyTargetType(targetTypeStr)
		trophies = append(trophies, t)
	}

	return trophies, rows.Err()
}

func (r *PostgresTrophyRepository) GetUserShowcase(ctx context.Context, userID string) ([]*entities.ShowcaseTrophy, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			t.id, t.user_id, t.name, t.description, t.icon, t.color, t.target_type, t.habit_id, t.required_streak_days,
			(ut.id IS NOT NULL) AS unlocked,
			ut.unlocked_at,
			t.created_at
		FROM trophies t
		LEFT JOIN user_trophies ut ON ut.trophy_id = t.id AND ut.user_id = $1
		WHERE t.user_id = $1
		ORDER BY unlocked DESC, t.created_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var showcase []*entities.ShowcaseTrophy
	for rows.Next() {
		item := &entities.ShowcaseTrophy{}
		var targetTypeStr string
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Description,
			&item.Icon,
			&item.Color,
			&targetTypeStr,
			&item.HabitID,
			&item.RequiredStreakDays,
			&item.Unlocked,
			&item.UnlockedAt,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.TargetType = entities.TrophyTargetType(targetTypeStr)
		showcase = append(showcase, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cargar condiciones para trofeos de sinergia
	for _, item := range showcase {
		if item.TargetType == entities.TargetSynergy {
			conds, err := r.GetConditionsByTrophyID(ctx, item.ID)
			if err == nil {
				item.Conditions = conds
			}
		}
	}

	return showcase, nil
}

func (r *PostgresTrophyRepository) UnlockTrophy(ctx context.Context, userID string, trophyID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO user_trophies (user_id, trophy_id, unlocked_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, trophy_id) DO NOTHING`

	cmdTag, err := r.db.Exec(ctx, query, userID, trophyID)
	if err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

func (r *PostgresTrophyRepository) IsTrophyUnlocked(ctx context.Context, userID string, trophyID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM user_trophies WHERE user_id = $1 AND trophy_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, trophyID).Scan(&exists)
	return exists, err
}

func (r *PostgresTrophyRepository) GetConditionsByTrophyID(ctx context.Context, trophyID string) ([]entities.TrophyCondition, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, trophy_id, habit_id, required_streak_days, created_at
		FROM trophy_conditions
		WHERE trophy_id = $1`

	rows, err := r.db.Query(ctx, query, trophyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conds []entities.TrophyCondition
	for rows.Next() {
		var c entities.TrophyCondition
		if err := rows.Scan(&c.ID, &c.TrophyID, &c.HabitID, &c.RequiredStreakDays, &c.CreatedAt); err != nil {
			return nil, err
		}
		conds = append(conds, c)
	}

	return conds, rows.Err()
}

func (r *PostgresTrophyRepository) FindSynergyTrophiesByHabit(ctx context.Context, habitID string) ([]*entities.Trophy, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT DISTINCT t.id, t.user_id, t.name, t.description, t.icon, t.color, t.target_type, t.habit_id, t.required_streak_days, t.created_at
		FROM trophies t
		INNER JOIN trophy_conditions tc ON tc.trophy_id = t.id
		WHERE tc.habit_id = $1 AND t.target_type = 'synergy'`

	rows, err := r.db.Query(ctx, query, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trophies []*entities.Trophy
	for rows.Next() {
		t := &entities.Trophy{}
		var targetTypeStr string
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Name,
			&t.Description,
			&t.Icon,
			&t.Color,
			&targetTypeStr,
			&t.HabitID,
			&t.RequiredStreakDays,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		t.TargetType = entities.TrophyTargetType(targetTypeStr)
		trophies = append(trophies, t)
	}

	return trophies, rows.Err()
}

func (r *PostgresTrophyRepository) FindActiveTrophiesByTarget(
	ctx context.Context,
	userID string,
	targetType entities.TrophyTargetType,
	habitID *string,
) ([]*entities.Trophy, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var query string
	var rows pgx.Rows
	var err error

	if habitID != nil {
		query = `
			SELECT t.id, t.user_id, t.name, t.description, t.icon, t.color, t.target_type, t.habit_id, t.required_streak_days, t.created_at
			FROM trophies t
			LEFT JOIN user_trophies ut ON ut.trophy_id = t.id AND ut.user_id = $1
			WHERE t.user_id = $1 AND t.target_type = $2 AND t.habit_id = $3 AND ut.id IS NULL`
		rows, err = r.db.Query(ctx, query, userID, string(targetType), *habitID)
	} else {
		query = `
			SELECT t.id, t.user_id, t.name, t.description, t.icon, t.color, t.target_type, t.habit_id, t.required_streak_days, t.created_at
			FROM trophies t
			LEFT JOIN user_trophies ut ON ut.trophy_id = t.id AND ut.user_id = $1
			WHERE t.user_id = $1 AND t.target_type = $2 AND ut.id IS NULL`
		rows, err = r.db.Query(ctx, query, userID, string(targetType))
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trophies []*entities.Trophy
	for rows.Next() {
		t := &entities.Trophy{}
		var targetTypeStr string
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Name,
			&t.Description,
			&t.Icon,
			&t.Color,
			&targetTypeStr,
			&t.HabitID,
			&t.RequiredStreakDays,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		t.TargetType = entities.TrophyTargetType(targetTypeStr)
		trophies = append(trophies, t)
	}

	return trophies, rows.Err()
}
