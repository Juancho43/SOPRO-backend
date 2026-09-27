package postgress

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SaveUser(user *entities.User) error {
	// Solo insertamos firebase_uid y email.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// id, streaks y timestamps se generan automáticamente en la BD por los valores DEFAULT[cite: 11].
	query := `
		INSERT INTO users (firebase_uid, email) 
		VALUES ($1, $2) 
		RETURNING id, current_streak, max_streak, created_at, updated_at`

	// Ejecutamos y recuperamos los valores autogenerados para actualizar la entidad en memoria
	err := r.db.QueryRow(ctx, query, user.UID, user.Email).Scan(
		&user.UID,
		&user.CurrentStreak,
		&user.MaxStreak,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}

func (r *PostgresUserRepository) GetUser(uid string) (*entities.User, error) {
	query := `
		SELECT id, firebase_uid, email, current_streak, max_streak, created_at, updated_at 
		FROM users 
		WHERE firebase_uid = $1`

	user := &entities.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Mapeamos todas las columnas de la tabla al struct[cite: 11]
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID,
		&user.Email,
		&user.CurrentStreak,
		&user.MaxStreak,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado") // Idealmente, devuelve un error personalizado de dominio
		}
		return nil, err
	}

	return user, nil
}
