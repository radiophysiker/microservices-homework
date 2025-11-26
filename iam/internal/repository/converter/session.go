package converter

import (
	serviceModel "github.com/radiophysiker/microservices-homework/iam/internal/model"
	repoModel "github.com/radiophysiker/microservices-homework/iam/internal/repository/model"
)

// ToRepoSession преобразует доменную модель сессии в модель repository слоя.
func ToRepoSession(session *serviceModel.Session) *repoModel.Session {
	if session == nil {
		return nil
	}

	return &repoModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
		IP:        session.IP,
		UserAgent: session.UserAgent,
		RevokedAt: session.RevokedAt,
	}
}

// ToServiceSession преобразует модель repository слоя в доменную модель сессии.
func ToServiceSession(session *repoModel.Session) *serviceModel.Session {
	if session == nil {
		return nil
	}

	return &serviceModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
		IP:        session.IP,
		UserAgent: session.UserAgent,
		RevokedAt: session.RevokedAt,
	}
}
