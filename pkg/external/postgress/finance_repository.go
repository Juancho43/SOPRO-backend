package postgress

import (
	"context"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresFinanceRepository implements repositories.FinanceRepository backed by PostgreSQL.
type PostgresFinanceRepository struct {
	db *pgxpool.Pool
}

// NewPostgresFinanceRepository constructs a new PostgresFinanceRepository instance.
func NewPostgresFinanceRepository(db *pgxpool.Pool) *PostgresFinanceRepository {
	return &PostgresFinanceRepository{db: db}
}

// CreateMonthlyGoal inserts a new monthly tactical goal into the database.
func (r *PostgresFinanceRepository) CreateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO monthly_goals (user_id, year, month, title, description, annual_goal_reference_id, status, achieved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		goal.UserID,
		goal.Year,
		goal.Month,
		goal.Title,
		goal.Description,
		goal.AnnualGoalReferenceID,
		string(goal.Status),
		goal.AchievedAt,
	).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
}

// GetMonthlyGoalByID retrieves a monthly goal by its ID and owner user ID.
func (r *PostgresFinanceRepository) GetMonthlyGoalByID(ctx context.Context, id string, userID string) (*entities.MonthlyGoal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, year, month, title, description, annual_goal_reference_id, status, achieved_at, created_at, updated_at
		FROM monthly_goals
		WHERE id = $1 AND user_id = $2`

	var g entities.MonthlyGoal
	var statusStr string
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&g.ID,
		&g.UserID,
		&g.Year,
		&g.Month,
		&g.Title,
		&g.Description,
		&g.AnnualGoalReferenceID,
		&statusStr,
		&g.AchievedAt,
		&g.CreatedAt,
		&g.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	g.Status = entities.MonthlyGoalStatus(statusStr)
	return &g, nil
}

// ListMonthlyGoalsByPeriod lists all monthly tactical goals for a user in a specific year and month.
func (r *PostgresFinanceRepository) ListMonthlyGoalsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.MonthlyGoal, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, year, month, title, description, annual_goal_reference_id, status, achieved_at, created_at, updated_at
		FROM monthly_goals
		WHERE user_id = $1 AND year = $2 AND month = $3
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*entities.MonthlyGoal
	for rows.Next() {
		var g entities.MonthlyGoal
		var statusStr string
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Year,
			&g.Month,
			&g.Title,
			&g.Description,
			&g.AnnualGoalReferenceID,
			&statusStr,
			&g.AchievedAt,
			&g.CreatedAt,
			&g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		g.Status = entities.MonthlyGoalStatus(statusStr)
		goals = append(goals, &g)
	}

	return goals, nil
}

// UpdateMonthlyGoal updates the title, description, and optional annual goal reference of a monthly goal.
func (r *PostgresFinanceRepository) UpdateMonthlyGoal(ctx context.Context, goal *entities.MonthlyGoal) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE monthly_goals
		SET title = $1, description = $2, annual_goal_reference_id = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		goal.Title,
		goal.Description,
		goal.AnnualGoalReferenceID,
		goal.ID,
		goal.UserID,
	).Scan(&goal.UpdatedAt)
}

// UpdateMonthlyGoalStatus updates the status and achieved_at timestamp of a monthly goal.
func (r *PostgresFinanceRepository) UpdateMonthlyGoalStatus(ctx context.Context, id string, userID string, status entities.MonthlyGoalStatus, achievedAt *time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE monthly_goals
		SET status = $1, achieved_at = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4`

	cmdTag, err := r.db.Exec(ctx, query, string(status), achievedAt, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("meta mensual no encontrada")
	}
	return nil
}

// DeleteMonthlyGoal removes a monthly goal owned by the specified user.
func (r *PostgresFinanceRepository) DeleteMonthlyGoal(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM monthly_goals WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("meta mensual no encontrada")
	}
	return nil
}

// CreateBudgetItem inserts a new budget item (income or expense) into the database.
func (r *PostgresFinanceRepository) CreateBudgetItem(ctx context.Context, item *entities.BudgetItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var categoryPtr *string
	if item.ExpenseCategory != nil {
		catStr := string(*item.ExpenseCategory)
		categoryPtr = &catStr
	}

	query := `
		INSERT INTO budget_items (user_id, year, month, type, expense_category, concept, amount, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		item.UserID,
		item.Year,
		item.Month,
		string(item.Type),
		categoryPtr,
		item.Concept,
		item.Amount,
		item.Currency,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

// GetBudgetItemByID retrieves a specific budget item by its ID and owner user ID.
func (r *PostgresFinanceRepository) GetBudgetItemByID(ctx context.Context, id string, userID string) (*entities.BudgetItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, year, month, type, expense_category, concept, amount, currency, created_at, updated_at
		FROM budget_items
		WHERE id = $1 AND user_id = $2`

	var item entities.BudgetItem
	var typeStr string
	var categoryPtr *string

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&item.ID,
		&item.UserID,
		&item.Year,
		&item.Month,
		&typeStr,
		&categoryPtr,
		&item.Concept,
		&item.Amount,
		&item.Currency,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	item.Type = entities.BudgetItemType(typeStr)
	if categoryPtr != nil {
		cat := entities.ExpenseCategory(*categoryPtr)
		item.ExpenseCategory = &cat
	}
	return &item, nil
}

// ListBudgetItemsByPeriod lists all budget items for a user in a given year and month.
func (r *PostgresFinanceRepository) ListBudgetItemsByPeriod(ctx context.Context, userID string, year int, month int) ([]*entities.BudgetItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, year, month, type, expense_category, concept, amount, currency, created_at, updated_at
		FROM budget_items
		WHERE user_id = $1 AND year = $2 AND month = $3
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entities.BudgetItem
	for rows.Next() {
		var item entities.BudgetItem
		var typeStr string
		var categoryPtr *string

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Year,
			&item.Month,
			&typeStr,
			&categoryPtr,
			&item.Concept,
			&item.Amount,
			&item.Currency,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		item.Type = entities.BudgetItemType(typeStr)
		if categoryPtr != nil {
			cat := entities.ExpenseCategory(*categoryPtr)
			item.ExpenseCategory = &cat
		}
		items = append(items, &item)
	}

	return items, nil
}

// UpdateBudgetItem updates an existing budget item in the database.
func (r *PostgresFinanceRepository) UpdateBudgetItem(ctx context.Context, item *entities.BudgetItem) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var categoryPtr *string
	if item.ExpenseCategory != nil {
		catStr := string(*item.ExpenseCategory)
		categoryPtr = &catStr
	}

	query := `
		UPDATE budget_items
		SET type = $1, expense_category = $2, concept = $3, amount = $4, currency = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query,
		string(item.Type),
		categoryPtr,
		item.Concept,
		item.Amount,
		item.Currency,
		item.ID,
		item.UserID,
	).Scan(&item.UpdatedAt)
}

// DeleteBudgetItem deletes a budget item owned by the specified user.
func (r *PostgresFinanceRepository) DeleteBudgetItem(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM budget_items WHERE id = $1 AND user_id = $2`
	cmdTag, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("partida presupuestaria no encontrada")
	}
	return nil
}
