-- +goose Up
-- SQL in section Connected with Analysis of Epic 3 (Challenges & Milestones)
CREATE TABLE IF NOT EXISTS challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,
    title VARCHAR(120) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reward TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_challenge_status CHECK (status IN ('ACTIVE', 'COMPLETED', 'FAILED', 'ARCHIVED', 'CANCELLED'))
);

CREATE INDEX IF NOT EXISTS idx_challenges_user_id ON challenges(user_id);
CREATE INDEX IF NOT EXISTS idx_challenges_user_status ON challenges(user_id, status);

CREATE TABLE IF NOT EXISTS milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    challenge_id UUID NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    due_date DATE NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_milestones_challenge_id ON milestones(challenge_id);
CREATE INDEX IF NOT EXISTS idx_milestones_due_date ON milestones(due_date);

-- +goose Down
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS challenges;
