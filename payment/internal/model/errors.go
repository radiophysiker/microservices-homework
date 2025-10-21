package model

import (
	"errors"
)

// ErrInvalidPaymentRequest - ошибка "некорректный запрос на оплату"
var ErrInvalidPaymentRequest = errors.New("invalid payment request")
