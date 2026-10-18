-- +goose Up
CREATE TABLE IF NOT EXISTS trophy_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trophy_id UUID NOT NULL,
    habit_id UUID NOT NULL,
    required_streak_days INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_trophy_conditions_trophy FOREIGN KEY (trophy_id) REFERENCES trophies(id) ON DELETE CASCADE,
    CONSTRAINT fk_trophy_conditions_habit FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
    CONSTRAINT unique_trophy_habit_condition UNIQUE (trophy_id, habit_id)
);

CREATE INDEX IF NOT EXISTS idx_trophy_conditions_trophy ON trophy_conditions(trophy_id);
CREATE INDEX IF NOT EXISTS idx_trophy_conditions_habit ON trophy_conditions(habit_id);

-- +goose Down
DROP TABLE IF EXISTS trophy_conditions;
