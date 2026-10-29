-- +goose Up
-- SQL in section Connected with Analysis of Epic 4 (Monthly Goals & Budget)
CREATE TABLE IF NOT EXISTS monthly_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    year INT NOT NULL,
    month INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    annual_goal_reference_id TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    achieved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_monthly_goal_month CHECK (month >= 1 AND month <= 12),
    CONSTRAINT check_monthly_goal_year CHECK (year >= 2000 AND year <= 2100),
    CONSTRAINT check_monthly_goal_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'ACHIEVED'))
);

CREATE INDEX IF NOT EXISTS idx_monthly_goals_user_id ON monthly_goals(user_id);
CREATE INDEX IF NOT EXISTS idx_monthly_goals_user_period ON monthly_goals(user_id, year, month);

CREATE TABLE IF NOT EXISTS budget_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    year INT NOT NULL,
    month INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    expense_category VARCHAR(20),
    concept VARCHAR(255) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_budget_item_month CHECK (month >= 1 AND month <= 12),
    CONSTRAINT check_budget_item_year CHECK (year >= 2000 AND year <= 2100),
    CONSTRAINT check_budget_item_type CHECK (type IN ('INCOME', 'EXPENSE')),
    CONSTRAINT check_budget_expense_category CHECK (
        (type = 'INCOME' AND expense_category IS NULL) OR
        (type = 'EXPENSE' AND expense_category IN ('FIXED', 'VARIABLE'))
    ),
    CONSTRAINT check_budget_amount_positive CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_budget_items_user_id ON budget_items(user_id);
CREATE INDEX IF NOT EXISTS idx_budget_items_user_period ON budget_items(user_id, year, month);

-- +goose Down
DROP TABLE IF EXISTS budget_items;
DROP TABLE IF EXISTS monthly_goals;
