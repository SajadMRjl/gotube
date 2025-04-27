package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RegisterHandlers(bot *Bot) {
	// Command handlers
	bot.RegisterHandler("start", startHandler)
	bot.RegisterHandler("help", helpHandler)
	bot.RegisterHandler("echo", echoHandler)

	// Default handler for non-command messages
	bot.RegisterHandler("default", defaultHandler)
}

func startHandler(ctx context.Context, bot *Bot, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, 
		"Welcome to the bot! I'm here to help. Use /help to see available commands.")
	msg.ReplyToMessageID = update.Message.MessageID

	_, err := bot.api.Send(msg)
	return err
}

func helpHandler(ctx context.Context, bot *Bot, update *tgbotapi.Update) error {
	helpText := `Available commands:
/start - Start the bot
/help - Show this help message
/echo <text> - Echo back the provided text`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	_, err := bot.api.Send(msg)
	return err
}

func echoHandler(ctx context.Context, bot *Bot, update *tgbotapi.Update) error {
	text := update.Message.CommandArguments()
	if text == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please provide text to echo after the /echo command")
		_, err := bot.api.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("You said: %s", text))
	_, err := bot.api.Send(msg)
	return err
}

func defaultHandler(ctx context.Context, bot *Bot, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, 
		"I don't understand that command. Try /help to see available commands.")
	_, err := bot.api.Send(msg)
	return err
}