-- +goose Up
CREATE TABLE habit_streaks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL,
    habit_name VARCHAR(255) NOT NULL,
    frequency VARCHAR(50) DEFAULT 'diario',
    current_streak INTEGER DEFAULT 0,
    max_streak INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_habit_streaks_user FOREIGN KEY (user_id) REFERENCES users(firebase_uid) ON DELETE CASCADE
);

CREATE INDEX idx_habit_streaks_user ON habit_streaks(user_id);

-- +goose Down
DROP TABLE IF EXISTS habit_streaks;