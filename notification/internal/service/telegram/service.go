package telegram

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"go.uber.org/zap"

	"github.com/radiophysiker/microservices-homework/notification/internal/client/http/telegram"
	"github.com/radiophysiker/microservices-homework/notification/internal/model"
	"github.com/radiophysiker/microservices-homework/platform/pkg/logger"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type Service interface {
	SendOrderPaidNotification(ctx context.Context, event *model.OrderPaid) error
	SendShipAssembledNotification(ctx context.Context, event *model.ShipAssembled) error
	HandleStartCommand(ctx context.Context, chatID string) error
}

type service struct {
	client        telegramClient
	chatID        string
	paidTmpl      *template.Template
	assembledTmpl *template.Template
}

type telegramClient interface {
	SendMessage(ctx context.Context, chatID, text string) error
	HandleStartCommand(ctx context.Context, chatID string) error
}

func NewService(client *telegram.Client, chatID string) (Service, error) {
	paidTemplateData, err := templatesFS.ReadFile("templates/paid_notification.tmpl")
	if err != nil {
		return nil, fmt.Errorf("read paid template: %w", err)
	}

	paidTmpl, err := template.New("paid").Parse(string(paidTemplateData))
	if err != nil {
		return nil, fmt.Errorf("parse paid template: %w", err)
	}

	assembledTemplateData, err := templatesFS.ReadFile("templates/assembled_notification.tmpl")
	if err != nil {
		return nil, fmt.Errorf("read assembled template: %w", err)
	}

	assembledTmpl, err := template.New("assembled").Parse(string(assembledTemplateData))
	if err != nil {
		return nil, fmt.Errorf("parse assembled template: %w", err)
	}

	return &service{
		client:        client,
		chatID:        chatID,
		paidTmpl:      paidTmpl,
		assembledTmpl: assembledTmpl,
	}, nil
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event *model.OrderPaid) error {
	var buf bytes.Buffer
	if err := s.paidTmpl.Execute(&buf, event); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	message := buf.String()

	if err := s.client.SendMessage(ctx, s.chatID, message); err != nil {
		logger.Error(ctx, "Failed to send OrderPaid notification",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return fmt.Errorf("send message: %w", err)
	}

	logger.Info(ctx, "OrderPaid notification sent",
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("chat_id", s.chatID),
	)

	return nil
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event *model.ShipAssembled) error {
	var buf bytes.Buffer
	if err := s.assembledTmpl.Execute(&buf, event); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	message := buf.String()

	if err := s.client.SendMessage(ctx, s.chatID, message); err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled notification",
			zap.Error(err),
			zap.String("order_uuid", event.OrderUUID.String()),
		)

		return fmt.Errorf("send message: %w", err)
	}

	logger.Info(ctx, "ShipAssembled notification sent",
		zap.String("order_uuid", event.OrderUUID.String()),
		zap.String("chat_id", s.chatID),
	)

	return nil
}

func (s *service) HandleStartCommand(ctx context.Context, chatID string) error {
	if err := s.client.HandleStartCommand(ctx, chatID); err != nil {
		logger.Error(ctx, "Failed to handle start command",
			zap.Error(err),
			zap.String("chat_id", chatID),
		)

		return fmt.Errorf("handle start command: %w", err)
	}

	logger.Info(ctx, "Start command handled successfully",
		zap.String("chat_id", chatID),
	)

	return nil
}
