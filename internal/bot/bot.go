package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	config   *Config
	logger   *zap.Logger
	handlers map[string]HandlerFunc
}

type HandlerFunc func(ctx context.Context, bot *Bot, update *tgbotapi.Update) error

func NewBot(config *Config, logger *zap.Logger) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = config.Telegram.Debug

	return &Bot{
		api:      botAPI,
		config:   config,
		logger:   logger,
		handlers: make(map[string]HandlerFunc),
	}, nil
}

func (b *Bot) RegisterHandler(command string, handler HandlerFunc) {
	b.handlers[command] = handler
}

func (b *Bot) StartPolling(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.config.Telegram.Timeout

	updates := b.api.GetUpdatesChan(u)

	b.logger.Info("Starting bot in polling mode")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			if err := b.handleUpdate(ctx, &update); err != nil {
				b.logger.Error("Error handling update", zap.Error(err))
			}
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	// Check if message is a command
	if update.Message.IsCommand() {
		handler, ok := b.handlers[update.Message.Command()]
		if ok {
			return handler(ctx, b, update)
		}
	}

	// Default handler for non-command messages
	if handler, ok := b.handlers["default"]; ok {
		return handler(ctx, b, update)
	}

	return nil
}