package main

import (
	"embed"

	"github.com/Juancho43/SOPRO-backend/pkg/core"
	"github.com/Juancho43/SOPRO-backend/pkg/src/injectors"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

// @title SOPRO_CORE
// @version 0.1
// @description SOPRO_CORE_ENGINE
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := core.LoadConfig()
	core.Migrate(cfg.DatabaseURL, embedMigrations)
	appCore := core.NewCore(cfg)
	injectors.InjectAll(appCore)
	appCore.Start()
}
