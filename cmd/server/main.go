package main

import (
	"DelayedNotifier/internal/config"
	"DelayedNotifier/internal/handlers"
	"DelayedNotifier/internal/repository"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.InitConsole()

	cfg := config.Load()
	zlog.Logger.Info().Str("port", cfg.HTTPPort).Msg("Starting notification service")

	db, err := dbpg.New(cfg.PostgreSQL.MasterDSN, cfg.PostgreSQL.SlaveDSNs, &dbpg.Options{
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("DB connection error")
	}
	defer db.Master.Close()

	repo := repository.NewPostgresRepository(db)

	notificationHandler := handlers.NewNotificationHandler(repo)

	r := ginext.New("debug")
	r.Use(ginext.Logger(), ginext.Recovery())

	api := r.Group("/api")
	{
		api.POST("/notify", notificationHandler.CreateNotification)
		api.GET("/notify/:id", notificationHandler.GetNotification)
		api.DELETE("/notify/:id", notificationHandler.CancelNotification)
	}

	zlog.Logger.Info().Str("port", cfg.HTTPPort).Msg("Server starting")
	err = r.Run(":" + cfg.HTTPPort)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("Start server error")
	}
}
