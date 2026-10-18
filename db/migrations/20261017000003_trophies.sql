-- +goose Up
CREATE TABLE IF NOT EXISTS trophies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100) DEFAULT 'military_tech',
    color VARCHAR(50) DEFAULT '#FFD700',
    target_type VARCHAR(50) NOT NULL, -- 'master_streak' | 'habit' | 'synergy'
    habit_id UUID NULL,
    required_streak_days INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_trophies_user FOREIGN KEY (user_id) REFERENCES users(firebase_uid) ON DELETE CASCADE,
    CONSTRAINT fk_trophies_habit FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_trophies_user ON trophies(user_id);

CREATE TABLE IF NOT EXISTS user_trophies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL,
    trophy_id UUID NOT NULL,
    unlocked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_trophies_user FOREIGN KEY (user_id) REFERENCES users(firebase_uid) ON DELETE CASCADE,
    CONSTRAINT fk_user_trophies_trophy FOREIGN KEY (trophy_id) REFERENCES trophies(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_trophy UNIQUE (user_id, trophy_id)
);

CREATE INDEX IF NOT EXISTS idx_user_trophies_user ON user_trophies(user_id);

-- +goose Down
DROP TABLE IF EXISTS user_trophies;
DROP TABLE IF EXISTS trophies;
