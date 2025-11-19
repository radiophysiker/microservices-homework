package model

import "time"

// Session — доменная сессия авторизации.
type Session struct {
	UUID     string
	UserUUID string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time

	IP        string
	UserAgent string
	RevokedAt *time.Time
}

// Credentials — доменные учётные данные
type Credentials struct {
	Login    string
	Password string
}
