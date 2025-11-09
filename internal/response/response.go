package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// Описываем структуру запроса, в которую будем парсить JSON из POST запроса
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// Описываем возвращаемый ответ для маршалинга
type Response struct {
	Alias  string `json:"alias,omitempty"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

// Конструируем ответ с ошибкой
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

// Конструируем ответ без ошибки
func OK(alias string) Response {
	return Response{
		Status: StatusOK,
		Alias : alias,
	}
}

//Конструируем читабельный ответ с ошибкой валидации

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, "\n"),
	}
}
