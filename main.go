package main

import (
	"fmt"
	"go-api/config"
	"go-api/controllers"
	"go-api/docs"
	"go-api/models"
	"go-api/routes"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type CLI struct {
	Port      int              `kong:"default='8080',help='Server port'"`
	Host      string           `kong:"default='localhost',help='Server host'"`
	DbPath    string           `kong:"default='app.db',help='SQLite database path'"`
	Debug     bool             `kong:"help='Enable debug mode'"`
	LogLevel  string           `kong:"default='info',enum='debug,info,warn,error',help='Log level (debug, info, warn, error)'"`
	LogFormat string           `kong:"default='text',enum='text,json',help='Log format (text, json)'"`
	Version   kong.VersionFlag `kong:"short='v',help='Show version'"`
}

// @title Your Project API
// @version 1.0
// @description This is a sample server for your project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("go-api"),
		kong.Description("A REST API server with Gin, GORM, and SQLite"),
		kong.Vars{
			"version": "1.0.0",
		},
	)

	// Setup structured logging
	logger := setupLogger(cli.LogLevel, cli.LogFormat)
	slog.SetDefault(logger)

	// Set Gin mode based on debug flag
	if cli.Debug {
		gin.SetMode(gin.DebugMode)
		slog.Debug("Debug mode enabled")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database with custom path
	database := config.InitDB(cli.DbPath, logger)

	// Auto migrate models
	err := database.AutoMigrate(&models.User{})
	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		ctx.FatalIfErrorf(err, "Failed to migrate database")
	}

	// Initialize Gin with custom logger middleware
	r := gin.New()
	//	r.Use(ginSlogMiddleware(logger))
	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())

	// Initialize controllers
	userController := controllers.NewUserController(database, logger)

	// Setup routes
	routes.SetupRoutes(r, userController)

	// Swagger endpoint
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = cli.Host + ":" + string(rune(cli.Port))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cli.Host, cli.Port)
	slog.Info("Starting server",
		"address", serverAddr,
		"debug", cli.Debug,
		"log_level", cli.LogLevel,
		"log_format", cli.LogFormat,
		"db_path", cli.DbPath,
	)

	if err := r.Run(serverAddr); err != nil {
		slog.Error("Failed to start server", "error", err, "address", serverAddr)
		ctx.FatalIfErrorf(err, "Failed to start server")
	}
}

// setupLogger configures slog with the specified level and format
func setupLogger(level, format string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
