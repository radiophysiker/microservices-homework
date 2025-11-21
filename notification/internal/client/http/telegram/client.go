package telegram

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

type Client struct {
	bot *bot.Bot
}

func NewClient(token string) (*Client, error) {
	b, err := bot.New(token)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	return &Client{
		bot: b,
	}, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return fmt.Errorf("parse chat_id: %w", err)
	}

	params := &bot.SendMessageParams{
		ChatID: chatIDInt,
		Text:   text,
	}

	message, err := c.bot.SendMessage(ctx, params)
	if err != nil {
		logger.Error(ctx, "Failed to send message to Telegram",
			zap.Error(err),
			zap.String("chat_id", chatID),
		)

		return fmt.Errorf("send message: %w", err)
	}

	logger.Info(ctx, "Message sent to Telegram",
		zap.String("chat_id", chatID),
		zap.Int("message_id", message.ID),
	)

	return nil
}

func (c *Client) HandleStartCommand(ctx context.Context, chatID string) error {
	welcomeMessage := "👋 Привет! Я бот для уведомлений о заказах космических кораблей.\n\n" +
		"Я буду отправлять вам уведомления о:\n" +
		"• Оплате заказов\n" +
		"• Завершении сборки кораблей\n\n" +
		"Ожидайте уведомлений!"

	return c.SendMessage(ctx, chatID, welcomeMessage)
}

// RegisterStartHandler регистрирует обработчик команды /start в боте
func (c *Client) RegisterStartHandler() {
	c.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		if err := c.HandleStartCommand(ctx, chatID); err != nil {
			logger.Error(ctx, "Failed to handle /start command in webhook",
				zap.Error(err),
				zap.Int64("chat_id", update.Message.Chat.ID),
			)
		}
	})
}

// WebhookHandler возвращает HTTP handler для обработки webhook от Telegram
func (c *Client) WebhookHandler() http.HandlerFunc {
	return c.bot.WebhookHandler()
}
