-- +goose Up
-- 1. Renombrar la tabla habit_streaks a habits
ALTER TABLE habit_streaks RENAME TO habits;

-- 2. Renombrar la columna habit_name a name para mantener consistencia
ALTER TABLE habits RENAME COLUMN habit_name TO name;

-- 3. Renombrar el índice existente
ALTER INDEX idx_habit_streaks_user RENAME TO idx_habits_user;

-- 4. Agregar las nuevas columnas requeridas para el dominio y archivado
ALTER TABLE habits 
    ADD COLUMN archived BOOLEAN DEFAULT FALSE,
    ADD COLUMN domain_id UUID NULL;

-- 5. Agregar la clave foránea y el índice para domain_id
ALTER TABLE habits 
    ADD CONSTRAINT fk_habits_domain FOREIGN KEY (domain_id) REFERENCES domains(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_habits_domain ON habits(domain_id);

-- 6. Crear la tabla habit_logs referenciando a la tabla recién renombrada
CREATE TABLE IF NOT EXISTS habit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL,
    user_id VARCHAR(128) NOT NULL,
    log_date DATE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_habit_logs_habit FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE,
    CONSTRAINT fk_habit_logs_user FOREIGN KEY (user_id) REFERENCES users(firebase_uid) ON DELETE CASCADE,
    CONSTRAINT unique_habit_daily_log UNIQUE (habit_id, log_date)
);

CREATE INDEX IF NOT EXISTS idx_habit_logs_query ON habit_logs(habit_id, log_date);
CREATE INDEX IF NOT EXISTS idx_habit_logs_user_date ON habit_logs(user_id, log_date);

-- +goose Down
-- 1. Eliminar la tabla habit_logs
DROP TABLE IF EXISTS habit_logs;

-- 2. Eliminar la clave foránea, índices y columnas agregadas a habits
ALTER TABLE habits DROP CONSTRAINT IF EXISTS fk_habits_domain;
DROP INDEX IF EXISTS idx_habits_domain;
ALTER TABLE habits 
    DROP COLUMN IF EXISTS domain_id,
    DROP COLUMN IF EXISTS archived;

-- 3. Revertir el nombre de la columna name a habit_name
ALTER TABLE habits RENAME COLUMN name TO habit_name;

-- 4. Revertir el nombre del índice
ALTER INDEX idx_habits_user RENAME TO idx_habit_streaks_user;

-- 5. Revertir el nombre de la tabla de vuelta a habit_streaks
ALTER TABLE habits RENAME TO habit_streaks;