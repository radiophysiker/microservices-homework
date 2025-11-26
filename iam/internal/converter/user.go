package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/radiophysiker/microservices-homework/iam/internal/model"
	commonpb "github.com/radiophysiker/microservices-homework/shared/pkg/proto/common/v1"
)

// ToProtoUser преобразует доменную модель User в protobuf User
func ToProtoUser(u *model.User) *commonpb.User {
	if u == nil {
		return nil
	}

	return &commonpb.User{
		Uuid:      u.UUID,
		Info:      ToProtoUserInfo(&u.Info),
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

// ToProtoUserInfo преобразует доменную модель UserInfo в protobuf UserInfo
func ToProtoUserInfo(info *model.UserInfo) *commonpb.UserInfo {
	if info == nil {
		return nil
	}

	return &commonpb.UserInfo{
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: ToProtoNotificationMethods(info.NotificationMethods),
	}
}

// ToProtoNotificationMethods преобразует слайс доменных NotificationMethod в protobuf
func ToProtoNotificationMethods(methods []model.NotificationMethod) []*commonpb.NotificationMethod {
	if methods == nil {
		return []*commonpb.NotificationMethod{}
	}

	protoMethods := make([]*commonpb.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		protoMethods = append(protoMethods, &commonpb.NotificationMethod{
			ProviderName: string(method.Provider),
			Target:       method.Target,
		})
	}

	return protoMethods
}

// FromProtoNotificationMethods преобразует слайс protobuf NotificationMethod в доменные модели
func FromProtoNotificationMethods(protoMethods []*commonpb.NotificationMethod) []model.NotificationMethod {
	if protoMethods == nil {
		return []model.NotificationMethod{}
	}

	methods := make([]model.NotificationMethod, 0, len(protoMethods))

	for _, protoMethod := range protoMethods {
		if protoMethod == nil {
			continue
		}

		methods = append(methods, model.NotificationMethod{
			Provider: model.NotificationProvider(protoMethod.ProviderName),
			Target:   protoMethod.Target,
			Primary:  false,
		})
	}

	return methods
}
