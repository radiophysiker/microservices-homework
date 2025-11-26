package session

import "fmt"

const (
	sessionKeyPattern      = "iam:sessions:%s"
	userSessionsKeyPattern = "iam:user-sessions:%s"
)

// sessionKey формирует ключ Redis для хранения сессии.
func sessionKey(sessionUUID string) string {
	return fmt.Sprintf(sessionKeyPattern, sessionUUID)
}

// userSessionsKey формирует ключ Redis для множества сессий пользователя.
func userSessionsKey(userUUID string) string {
	return fmt.Sprintf(userSessionsKeyPattern, userUUID)
}
