package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SajadMRjl/gotube/internal/bot"
	"github.com/SajadMRjl/gotube/internal/storage"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	config, err := bot.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Initialize logger
	logger, err := bot.NewLogger(config.Logging.Level, config.Logging.Development)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("failed to sync logger: %v\n", err)
		}
	}()

	// Initialize database
	db, err := initDatabase(config, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer func() {
		sqlDB, err := db.DB.DB()
		if err != nil {
			logger.Error("Failed to get SQL DB for closing", zap.Error(err))
			return
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	// Create bot instance with dependencies
	botInstance, err := bot.NewBot(config, logger, db)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}

	// Register handlers
	bot.RegisterHandlers(botInstance)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupGracefulShutdown(ctx, cancel, logger)

	// Start the bot
	logger.Info("Starting bot...",
		zap.String("mode", "polling"),
		zap.Bool("debug", config.Telegram.Debug),
	)

	if err := botInstance.StartPolling(ctx); err != nil {
		logger.Fatal("Bot stopped with error", zap.Error(err))
	}
}

func initDatabase(config *bot.Config, logger *zap.Logger) (*storage.GormStorage, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
		config.Database.Timezone,
	)

	logger.Info("Connecting to database...",
		zap.String("host", config.Database.Host),
		zap.Int("port", config.Database.Port),
		zap.String("dbname", config.Database.DBName),
	)

	db, err := storage.NewGormStorage(connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)

	connMaxLifetime, err := time.ParseDuration(config.Database.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("invalid connection max lifetime: %w", err)
	}
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	logger.Info("Database connection established")

	// Run migrations using the models from the storage package
	logger.Info("Running database migrations...")
	if err := db.DB.AutoMigrate(
		&storage.User{},
		&storage.Track{},
	); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	logger.Info("Migrations completed")

	return db, nil
}

func setupGracefulShutdown(ctx context.Context, cancel context.CancelFunc, logger *zap.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal",
			zap.String("signal", sig.String()))

		// Cancel the root context
		cancel()

		// Give some time for cleanup
		select {
		case <-time.After(5 * time.Second):
			logger.Warn("Cleanup timeout exceeded, forcing shutdown")
		case <-ctx.Done():
			logger.Info("Cleanup completed successfully")
		}

		os.Exit(0)
	}()
}
