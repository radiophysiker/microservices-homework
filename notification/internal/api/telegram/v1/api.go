package v1

import (
	"net/http"

	"github.com/radiophysiker/microservices-homework/notification/internal/client/http/telegram"
)

type API struct {
	telegramClient *telegram.Client
}

func NewAPI(telegramClient *telegram.Client) *API {
	return &API{
		telegramClient: telegramClient,
	}
}

func (a *API) RegisterRoutes(mux *http.ServeMux) {
	// Используем webhook handler из библиотеки go-telegram/bot
	mux.HandleFunc("/webhook", a.telegramClient.WebhookHandler())
}
