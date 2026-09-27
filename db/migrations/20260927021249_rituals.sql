-- +goose Up
CREATE TABLE daily_rituals (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
ritual_date DATE NOT NULL,
gratitude TEXT NOT NULL,
goals JSONB NOT NULL,
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
-- Regla inquebrantable: Solo un ritual por día calendario
CONSTRAINT unique_user_daily_ritual UNIQUE (user_id, ritual_date)
);

CREATE INDEX idx_daily_rituals_user_date ON daily_rituals(user_id, ritual_date);

-- +goose Down
DROP TABLE IF EXISTS daily_rituals;
