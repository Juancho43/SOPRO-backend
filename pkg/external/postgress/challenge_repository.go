package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresChallengeRepository struct {
	db *pgxpool.Pool
}

func NewPostgresChallengeRepository(db *pgxpool.Pool) *PostgresChallengeRepository {
	return &PostgresChallengeRepository{db: db}
}

// CreateChallenge guarda un nuevo reto en la base de datos.
func (r *PostgresChallengeRepository) CreateChallenge(ctx context.Context, challenge *entities.Challenge) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO challenges (user_id, title, description, start_date, end_date, reward, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		challenge.UserID,
		challenge.Title,
		challenge.Description,
		challenge.StartDate,
		challenge.EndDate,
		challenge.Reward,
		string(challenge.Status),
	).Scan(&challenge.ID, &challenge.CreatedAt, &challenge.UpdatedAt)

	return err
}

// GetChallengeByID busca un reto específico de un usuario.
func (r *PostgresChallengeRepository) GetChallengeByID(ctx context.Context, id string, userID string) (*entities.Challenge, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, title, description, TO_CHAR(start_date, 'YYYY-MM-DD'), TO_CHAR(end_date, 'YYYY-MM-DD'), reward, status, created_at, updated_at
		FROM challenges
		WHERE id = $1 AND user_id = $2`

	var c entities.Challenge
	var statusStr string
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&c.ID,
		&c.UserID,
		&c.Title,
		&c.Description,
		&c.StartDate,
		&c.EndDate,
		&c.Reward,
		&statusStr,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	c.Status = entities.ChallengeStatus(statusStr)
	return &c, nil
}

// ListChallengesByUserID lista los retos del usuario con filtro de estado opcional.
func (r *PostgresChallengeRepository) ListChallengesByUserID(ctx context.Context, userID string, status *entities.ChallengeStatus) ([]*entities.Challenge, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var query string
	var rows pgx.Rows
	var err error

	if status != nil {
		query = `
			SELECT id, user_id, title, description, TO_CHAR(start_date, 'YYYY-MM-DD'), TO_CHAR(end_date, 'YYYY-MM-DD'), reward, status, created_at, updated_at
			FROM challenges
			WHERE user_id = $1 AND status = $2
			ORDER BY created_at DESC`
		rows, err = r.db.Query(ctx, query, userID, string(*status))
	} else {
		query = `
			SELECT id, user_id, title, description, TO_CHAR(start_date, 'YYYY-MM-DD'), TO_CHAR(end_date, 'YYYY-MM-DD'), reward, status, created_at, updated_at
			FROM challenges
			WHERE user_id = $1
			ORDER BY created_at DESC`
		rows, err = r.db.Query(ctx, query, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entities.Challenge
	for rows.Next() {
		var c entities.Challenge
		var statusStr string
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Title,
			&c.Description,
			&c.StartDate,
			&c.EndDate,
			&c.Reward,
			&statusStr,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.Status = entities.ChallengeStatus(statusStr)
		result = append(result, &c)
	}

	return result, nil
}

// UpdateChallenge actualiza título, descripción y recompensa del reto.
func (r *PostgresChallengeRepository) UpdateChallenge(ctx context.Context, challenge *entities.Challenge) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE challenges
		SET title = $1, description = $2, reward = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		challenge.Title,
		challenge.Description,
		challenge.Reward,
		challenge.ID,
		challenge.UserID,
	).Scan(&challenge.UpdatedAt)
}

// UpdateChallengeStatus cambia el estado del reto.
func (r *PostgresChallengeRepository) UpdateChallengeStatus(ctx context.Context, id string, userID string, status entities.ChallengeStatus) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE challenges
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3`

	cmdTag, err := r.db.Exec(ctx, query, string(status), id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("reto no encontrado para actualizar estado")
	}
	return nil
}

// CreateMilestone guarda un nuevo hito.
func (r *PostgresChallengeRepository) CreateMilestone(ctx context.Context, milestone *entities.Milestone) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO milestones (challenge_id, title, due_date, is_completed, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		milestone.ChallengeID,
		milestone.Title,
		milestone.DueDate,
		milestone.IsCompleted,
		milestone.OrderIndex,
	).Scan(&milestone.ID, &milestone.CreatedAt, &milestone.UpdatedAt)
}

// GetMilestoneByID obtiene un hito por su id.
func (r *PostgresChallengeRepository) GetMilestoneByID(ctx context.Context, id string) (*entities.Milestone, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, challenge_id, title, TO_CHAR(due_date, 'YYYY-MM-DD'), is_completed, completed_at, order_index, created_at, updated_at
		FROM milestones
		WHERE id = $1`

	var m entities.Milestone
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.ChallengeID,
		&m.Title,
		&m.DueDate,
		&m.IsCompleted,
		&m.CompletedAt,
		&m.OrderIndex,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// ListMilestonesByChallengeID obtiene todos los hitos ordenados cronológicamente por fecha límite.
func (r *PostgresChallengeRepository) ListMilestonesByChallengeID(ctx context.Context, challengeID string) ([]*entities.Milestone, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, challenge_id, title, TO_CHAR(due_date, 'YYYY-MM-DD'), is_completed, completed_at, order_index, created_at, updated_at
		FROM milestones
		WHERE challenge_id = $1
		ORDER BY due_date ASC, order_index ASC`

	rows, err := r.db.Query(ctx, query, challengeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []*entities.Milestone
	for rows.Next() {
		var m entities.Milestone
		if err := rows.Scan(
			&m.ID,
			&m.ChallengeID,
			&m.Title,
			&m.DueDate,
			&m.IsCompleted,
			&m.CompletedAt,
			&m.OrderIndex,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		milestones = append(milestones, &m)
	}

	return milestones, nil
}

// UpdateMilestone actualiza datos de un hito.
func (r *PostgresChallengeRepository) UpdateMilestone(ctx context.Context, milestone *entities.Milestone) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE milestones
		SET title = $1, due_date = $2, is_completed = $3, completed_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND challenge_id = $6
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		milestone.Title,
		milestone.DueDate,
		milestone.IsCompleted,
		milestone.CompletedAt,
		milestone.ID,
		milestone.ChallengeID,
	).Scan(&milestone.UpdatedAt)
}

// DeleteMilestone elimina un hito de un reto.
func (r *PostgresChallengeRepository) DeleteMilestone(ctx context.Context, id string, challengeID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM milestones WHERE id = $1 AND challenge_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, challengeID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("hito no encontrado")
	}
	return nil
}

// ListUpcomingMilestones busca hitos pendientes dentro de los próximos N días para retos activos del usuario.
func (r *PostgresChallengeRepository) ListUpcomingMilestones(ctx context.Context, userID string, maxDaysAhead int, today string) ([]*entities.UpcomingDeadline, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			m.id, 
			m.challenge_id, 
			c.title AS challenge_title, 
			m.title, 
			TO_CHAR(m.due_date, 'YYYY-MM-DD') AS due_date,
			(m.due_date - $2::date) AS days_remaining
		FROM milestones m
		INNER JOIN challenges c ON m.challenge_id = c.id
		WHERE c.user_id = $1 
		  AND c.status = 'ACTIVE' 
		  AND m.is_completed = FALSE
		  AND m.due_date <= ($2::date + $3 * INTERVAL '1 day')
		ORDER BY m.due_date ASC`

	rows, err := r.db.Query(ctx, query, userID, today, maxDaysAhead)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entities.UpcomingDeadline
	for rows.Next() {
		var item entities.UpcomingDeadline
		if err := rows.Scan(
			&item.MilestoneID,
			&item.ChallengeID,
			&item.ChallengeTitle,
			&item.Title,
			&item.DueDate,
			&item.DaysRemaining,
		); err != nil {
			return nil, err
		}

		if item.DaysRemaining < 0 {
			item.UrgencyLevel = "overdue"
		} else if item.DaysRemaining <= 1 {
			item.UrgencyLevel = "critical"
		} else if item.DaysRemaining <= 3 {
			item.UrgencyLevel = "warning"
		} else {
			item.UrgencyLevel = "normal"
		}

		list = append(list, &item)
	}

	return list, nil
}
