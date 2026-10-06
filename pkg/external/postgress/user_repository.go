package postgress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Juancho43/SOPRO-backend/pkg/internal/model/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SaveUser(user *entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
		INSERT INTO users (firebase_uid, email) 
		VALUES ($1, $2) 
		RETURNING firebase_uid, email, current_streak, max_streak, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, user.UID, user.Email).Scan(
		&user.UID,
		&user.Email,
		&user.CurrentStreak,
		&user.MaxStreak,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return err
}

func (r *PostgresUserRepository) GetUser(uid string) (*entities.User, error) {
	query := `
		SELECT firebase_uid, email, current_streak, max_streak, created_at, updated_at 
		FROM users 
		WHERE firebase_uid = $1`
	fmt.Print("id a buscar", uid)

	user := &entities.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID,
		&user.Email,
		&user.CurrentStreak,
		&user.MaxStreak,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return user, nil
}
