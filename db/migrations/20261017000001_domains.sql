-- +goose Up
CREATE TABLE IF NOT EXISTS domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20) NOT NULL DEFAULT '#3f51b5',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_domains_user FOREIGN KEY (user_id) REFERENCES users(firebase_uid) ON DELETE CASCADE,
    CONSTRAINT unique_user_domain_name UNIQUE (user_id, name)
);

CREATE INDEX IF NOT EXISTS idx_domains_user ON domains(user_id);

-- +goose Down
DROP TABLE IF EXISTS domains;
