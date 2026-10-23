package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDomainRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDomainRepository(db *pgxpool.Pool) *PostgresDomainRepository {
	return &PostgresDomainRepository{db: db}
}

func (r *PostgresDomainRepository) Create(ctx context.Context, domain *entities.Domain) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO domains (user_id, name, color) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, domain.UserID, domain.Name, domain.Color).Scan(
		&domain.ID,
		&domain.CreatedAt,
	)
}

func (r *PostgresDomainRepository) GetByID(ctx context.Context, id string) (*entities.Domain, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, name, color, created_at 
		FROM domains 
		WHERE id = $1`

	domain := &entities.Domain{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&domain.ID,
		&domain.UserID,
		&domain.Name,
		&domain.Color,
		&domain.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return domain, nil
}

func (r *PostgresDomainRepository) ListByUserID(ctx context.Context, userID string) ([]*entities.DomainWithHabitCount, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT d.id, d.user_id, d.name, d.color, d.created_at, COUNT(h.id) AS habit_count
		FROM domains d
		LEFT JOIN habits h ON h.domain_id = d.id AND h.archived = FALSE
		WHERE d.user_id = $1
		GROUP BY d.id, d.user_id, d.name, d.color, d.created_at
		ORDER BY d.created_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []*entities.DomainWithHabitCount
	for rows.Next() {
		item := &entities.DomainWithHabitCount{}
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Name,
			&item.Color,
			&item.CreatedAt,
			&item.HabitCount,
		); err != nil {
			return nil, err
		}
		domains = append(domains, item)
	}

	return domains, rows.Err()
}

func (r *PostgresDomainRepository) Update(ctx context.Context, domain *entities.Domain) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE domains 
		SET name = $1, color = $2 
		WHERE id = $3 AND user_id = $4`

	cmdTag, err := r.db.Exec(ctx, query, domain.Name, domain.Color, domain.ID, domain.UserID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("dominio no encontrado o no pertenece al usuario")
	}

	return nil
}

func (r *PostgresDomainRepository) Delete(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM domains WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("dominio no encontrado o no pertenece al usuario")
	}

	return nil
}

func (r *PostgresDomainRepository) GetWeeklyEffortByDomain(ctx context.Context, userID string, fromDate time.Time) ([]*entities.DomainEffortPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fromDateStr := fromDate.Format("2006-01-02")
	query := `
		SELECT d.id, d.name, d.color, COUNT(hl.id) AS completed_count
		FROM domains d
		LEFT JOIN habits h ON h.domain_id = d.id AND h.user_id = $1 AND h.archived = FALSE
		LEFT JOIN habit_logs hl ON hl.habit_id = h.id AND hl.log_date >= $2
		WHERE d.user_id = $1
		GROUP BY d.id, d.name, d.color
		ORDER BY completed_count DESC`

	rows, err := r.db.Query(ctx, query, userID, fromDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*entities.DomainEffortPoint
	for rows.Next() {
		point := &entities.DomainEffortPoint{}
		if err := rows.Scan(
			&point.DomainID,
			&point.DomainName,
			&point.Color,
			&point.CompletedCount,
		); err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, rows.Err()
}
