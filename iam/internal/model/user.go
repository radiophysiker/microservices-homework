package model

import "time"

type NotificationProvider string

const (
	NotificationProviderEmail    NotificationProvider = "email"
	NotificationProviderTelegram NotificationProvider = "telegram"
	NotificationProviderPush     NotificationProvider = "push"
)

// NotificationMethod - метод уведомления пользователя
type NotificationMethod struct {
	Provider NotificationProvider
	Target   string
	Primary  bool
}

// UserInfo - основная информация пользователя
type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

// User - агрегированные данные пользователя
type User struct {
	UUID         string
	Info         UserInfo
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}
