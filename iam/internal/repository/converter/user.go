package converter

import (
	"encoding/json"

	serviceModel "github.com/radiophysiker/microservices-homework/iam/internal/model"
	repoModel "github.com/radiophysiker/microservices-homework/iam/internal/repository/model"
)

// ToRepoUser преобразует доменную модель пользователя в модель repository слоя.
func ToRepoUser(user *serviceModel.User) (*repoModel.User, error) {
	if user == nil {
		return nil, nil
	}

	notificationJSON, err := json.Marshal(user.Info.NotificationMethods)
	if err != nil {
		return nil, err
	}

	return &repoModel.User{
		UUID:                user.UUID,
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		PasswordHash:        user.PasswordHash,
		NotificationMethods: notificationJSON,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}, nil
}

// ToServiceUser преобразует модель repository слоя в доменную модель пользователя.
func ToServiceUser(user *repoModel.User) (*serviceModel.User, error) {
	if user == nil {
		return nil, nil
	}

	var notificationMethods []serviceModel.NotificationMethod
	if len(user.NotificationMethods) > 0 {
		if err := json.Unmarshal(user.NotificationMethods, &notificationMethods); err != nil {
			return nil, err
		}
	}

	return &serviceModel.User{
		UUID: user.UUID,
		Info: serviceModel.UserInfo{
			Login:               user.Login,
			Email:               user.Email,
			NotificationMethods: notificationMethods,
		},
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
