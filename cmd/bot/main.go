package main

import (
	"attendance-bot/internal/attendance"
	"attendance-bot/internal/bot"
	"attendance-bot/internal/config"
	"attendance-bot/internal/database"
	"attendance-bot/internal/reports"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Configuration loaded", "environment", cfg.Environment)

	// Initialize database
	db, err := database.NewSQLiteDB(cfg.DatabasePath)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Database initialized", "path", cfg.DatabasePath)

	// Initialize repository
	repo := database.NewRepository(db)

	// Initialize attendance service
	attendanceService := attendance.NewService(repo, cfg.TOTPSecret)

	// Initialize CSV generator
	csvGenerator := reports.NewCSVGenerator("temp")

	// Initialize bot
	botInstance := bot.NewBot(cfg.BotToken, attendanceService, csvGenerator, logger)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start bot in a goroutine
	go func() {
		if err := botInstance.Start(); err != nil {
			logger.Error("Bot error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutting down gracefully...")
}
