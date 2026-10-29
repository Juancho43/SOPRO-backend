package entities

import "time"

// MonthlyGoalStatus represents the possible lifecycle states of a monthly tactical goal (HU 4.1).
type MonthlyGoalStatus string

const (
	// MonthlyGoalStatusPending indicates the monthly goal has been created and awaits execution.
	MonthlyGoalStatusPending MonthlyGoalStatus = "PENDING"
	// MonthlyGoalStatusInProgress indicates the monthly goal is actively being worked on.
	MonthlyGoalStatusInProgress MonthlyGoalStatus = "IN_PROGRESS"
	// MonthlyGoalStatusAchieved indicates the monthly goal has been successfully completed.
	MonthlyGoalStatusAchieved MonthlyGoalStatus = "ACHIEVED"
)

// MonthlyGoal represents a tactical monthly objective optionally linked to an annual ritual goal (HU 4.1).
type MonthlyGoal struct {
	ID                    string            `json:"id"`
	UserID                string            `json:"user_id"`
	Year                  int               `json:"year"`
	Month                 int               `json:"month"`
	Title                 string            `json:"title"`
	Description           string            `json:"description"`
	AnnualGoalReferenceID *string           `json:"annual_goal_reference_id,omitempty"`
	Status                MonthlyGoalStatus `json:"status"`
	AchievedAt            *time.Time        `json:"achieved_at,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

// BudgetItemType defines whether a budget line item is an income or an expense (HU 4.2).
type BudgetItemType string

const (
	// BudgetItemTypeIncome represents a projected or realized monthly income stream.
	BudgetItemTypeIncome BudgetItemType = "INCOME"
	// BudgetItemTypeExpense represents a monthly expense item.
	BudgetItemTypeExpense BudgetItemType = "EXPENSE"
)

// ExpenseCategory classifies an expense into fixed (structural) or variable (operational/lifestyle) (HU 4.2).
type ExpenseCategory string

const (
	// ExpenseCategoryFixed represents recurring structural expenses that form the break-even baseline.
	ExpenseCategoryFixed ExpenseCategory = "FIXED"
	// ExpenseCategoryVariable represents flexible operational or lifestyle expenses.
	ExpenseCategoryVariable ExpenseCategory = "VARIABLE"
)

// BudgetItem represents a single monthly budget entry for income or expense tracking (HU 4.2).
type BudgetItem struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Year            int              `json:"year"`
	Month           int              `json:"month"`
	Type            BudgetItemType   `json:"type"`
	ExpenseCategory *ExpenseCategory `json:"expense_category,omitempty"`
	Concept         string           `json:"concept"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// BudgetSummary encapsulates the cash flow, break-even point, and solvency metrics for a given month (HU 4.2).
type BudgetSummary struct {
	Year                  int     `json:"year"`
	Month                 int     `json:"month"`
	TotalIncome           float64 `json:"total_income"`
	TotalFixedExpenses    float64 `json:"total_fixed_expenses"` // Break-even baseline
	TotalVariableExpenses float64 `json:"total_variable_expenses"`
	TotalExpenses         float64 `json:"total_expenses"`
	NetBalance            float64 `json:"net_balance"`
	IsSolvent             bool    `json:"is_solvent"`
	DeficitGap            float64 `json:"deficit_gap"`    // Monetary gap to cover total expenses (0 if solvent)
	BreakEvenGap          float64 `json:"break_even_gap"` // Monetary gap to cover fixed expenses (0 if income >= fixed expenses)
}

// BudgetOverview aggregates the list of monthly budget items and the calculated financial summary.
type BudgetOverview struct {
	BudgetSummary
	Items   []*BudgetItem `json:"items"`
	Summary BudgetSummary `json:"summary"`
}
