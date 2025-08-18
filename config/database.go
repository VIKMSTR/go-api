package config

import (
	"log/slog"

	"github.com/glebarez/sqlite" // slower but portable sqlite driver, that does not need CGO. In case of high traffic, consider using non portable CGO one
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(dbPath string, log *slog.Logger) *gorm.DB {
	// Configure GORM logger to use slog
	gormLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Error("Failed to connect to database", "error", err, "path", dbPath)
		panic(err)
	}

	log.Info("Database connected successfully", "path", dbPath)
	return db
}
