package apperror

import "fmt"

type AppError struct {
	Code    string // код ошибки
	Message string // сообщение об ошибке
	Err     error  // ошибка
}

// Конструкторы на два кейса:
// -- 1. Простая ошибка
// -- 2. Ошибка с вложенной ошибкой
func New(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
