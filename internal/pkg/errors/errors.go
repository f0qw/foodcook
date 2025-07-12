package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// 预定义错误
var (
	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "资源未找到",
	}

	ErrBadRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "请求参数错误",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问",
	}

	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "禁止访问",
	}

	ErrInternalServer = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
	}

	ErrConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "资源冲突",
	}
)

// 创建错误的辅助函数
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func WrapError(err error, message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
