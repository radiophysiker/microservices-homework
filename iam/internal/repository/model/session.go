package model

import "time"

// Session представляет модель сессии в Redis.
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
