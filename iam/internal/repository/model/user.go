package model

import "time"

// User представляет модель пользователя на уровне repository (PostgreSQL).
type User struct {
	UUID                string
	Login               string
	Email               string
	PasswordHash        string
	NotificationMethods []byte

	CreatedAt time.Time
	UpdatedAt time.Time
}
