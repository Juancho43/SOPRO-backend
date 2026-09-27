package main

import (
	"embed"
	"fmt"

	"github.com/Juancho43/SOPRO-backend/app/core"
	"github.com/Juancho43/SOPRO-backend/app/src/injectors"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

// @title SOPRO_CORE
// @version 0.1
// @description SOPRO_CORE_ENGINE
// @BasePath /
func main() {
	cfg := core.LoadConfig()
	core.Migrate(cfg.DatabaseURL, embedMigrations)
	appCore := core.NewCore(cfg)
	injectors.InjectAll(appCore)
	fmt.Print(cfg)
	appCore.Start()
	fmt.Print("HI")
}
