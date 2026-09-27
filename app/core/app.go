package core

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Core struct {
	Router   *gin.Engine
	Config   *Config
	Logger   *slog.Logger
	Database *pgxpool.Pool
}

// NewCore inicializa el núcleo del sistema, instanciando el motor HTTP.
func NewCore(cfg *Config) *Core {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(setupCORS())
	return &Core{
		Router:   router,
		Config:   cfg,
		Logger:   NewLogger(cfg),
		Database: NewDatabasePool(cfg.DatabaseURL),
	}
}

// Start pone en marcha tu servidor. La acción es lo que produce resultados.
func (c *Core) Start() {
	c.Logger.Debug("Iniciando el servidor en el puerto" + c.Config.Port)
	if err := c.Router.Run(":" + c.Config.Port); err != nil {
		c.Logger.Error("Error al arrancar el servidor: " + err.Error())
	}
}
func setupCORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	// Define con exactitud quién tiene acceso a tus recursos
	corsConfig.AllowOrigins = []string{"http://localhost:4200", "https://sopro.juancetodev.com.ar"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	return cors.New(corsConfig)
}
