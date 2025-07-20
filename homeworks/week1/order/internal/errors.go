package internal

import (
	orderv1 "github.com/radiophysiker/microservices-homework/week1/shared/pkg/openapi/order/v1"
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Type    string
	Message string
	Code    int
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewOrderNotFoundError создает ошибку "заказ не найден"
func NewOrderNotFoundError(message string) *orderv1.NotFoundError {
	return &orderv1.NotFoundError{
		Error:   orderv1.NotFoundErrorErrorNotFound,
		Message: message,
	}
}

// NewOrderBadRequestError создает ошибку "неверный запрос"
func NewOrderBadRequestError(message string) *orderv1.BadRequestError {
	return &orderv1.BadRequestError{
		Error:   orderv1.BadRequestErrorErrorBadRequest,
		Message: message,
	}
}

// NewOrderConflictError создает ошибку "конфликт"
func NewOrderConflictError(message string) *orderv1.ConflictError {
	return &orderv1.ConflictError{
		Error:   orderv1.ConflictErrorErrorConflict,
		Message: message,
	}
}

// NewOrderInternalError создает ошибку "внутренняя ошибка сервера"
func NewOrderInternalError(message string) *orderv1.InternalServerError {
	return &orderv1.InternalServerError{
		Error:   orderv1.InternalServerErrorErrorInternalServerError,
		Message: message,
	}
}

// NewOrderGenericError создает общую ошибку
func NewOrderGenericError(statusCode int, message string) *orderv1.GenericErrorStatusCode {
	return &orderv1.GenericErrorStatusCode{
		StatusCode: statusCode,
		Response: orderv1.GenericError{
			Error:   "error",
			Message: message,
		},
	}
}

// PayOrder specific error functions
// NewPayOrderNotFoundError создает ошибку "заказ не найден" для PayOrder
func NewPayOrderNotFoundError(message string) orderv1.PayOrderRes {
	return &orderv1.NotFoundError{
		Error:   orderv1.NotFoundErrorErrorNotFound,
		Message: message,
	}
}

// NewPayOrderConflictError создает ошибку "конфликт" для PayOrder
func NewPayOrderConflictError(message string) orderv1.PayOrderRes {
	return &orderv1.ConflictError{
		Error:   orderv1.ConflictErrorErrorConflict,
		Message: message,
	}
}

// NewPayOrderInternalError создает ошибку "внутренняя ошибка сервера" для PayOrder
func NewPayOrderInternalError(message string) orderv1.PayOrderRes {
	return &orderv1.InternalServerError{
		Error:   orderv1.InternalServerErrorErrorInternalServerError,
		Message: message,
	}
}
