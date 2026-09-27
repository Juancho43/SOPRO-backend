package core

import (
	"context"
	"database/sql"
	"embed"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewDatabasePool(url string) *pgxpool.Pool {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalf("Fallo en la conexión: %v\n", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("El ping a Supabase falló: %v\n", err)
	}
	return pool
}

func Migrate(url string, migrationsFS embed.FS) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Error conectando a la BD: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error configurando dialecto: %v", err)
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	log.Println("Migraciones aplicadas con éxito en Supabase.")
}
