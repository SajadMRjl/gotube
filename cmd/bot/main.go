package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sajadMRjl/gotube/internal/bot"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	config, err := bot.LoadConfig("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	logger, err := bot.NewLogger(config.Logging.Level, config.Logging.Development)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Create bot instance
	botInstance, err := bot.NewBot(config, logger)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}

	// Register handlers
	bot.RegisterHandlers(botInstance)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down gracefully...")
		cancel()

		// Additional cleanup can go here
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// Start the bot in polling mode
	logger.Info("Starting bot...")
	if err := botInstance.StartPolling(ctx); err != nil {
		logger.Fatal("Bot stopped with error", zap.Error(err))
	}
}
