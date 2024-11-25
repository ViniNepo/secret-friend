package handler

import (
	"errors"
	"log"
	"net/http"
)

type ErrorHandlerConfig struct {
	Err        error
	HTTPStatus int
	ErrorCode  string
}

var (
	ErrRequestBodyIsInvalid  = NewError("request body is not a valid JSON")
	ErrInvalidQueryParameter = NewError("invalid query parameter")
	ErrSendEmail             = NewError("cannot send email")
)

var ErrorHandlerConfigs = []ErrorHandlerConfig{
	{
		Err:        ErrRequestBodyIsInvalid,
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  "json_decode_error",
	},
	{
		Err:        ErrInvalidQueryParameter,
		HTTPStatus: http.StatusBadRequest,
		ErrorCode:  "invalid_query_parameter",
	},
	{
		Err:        ErrSendEmail,
		HTTPStatus: http.StatusInternalServerError,
		ErrorCode:  "error_send_email",
	},
}

func NewError(message string) error {
	return &CustomError{Message: message}
}

// CustomError define a estrutura do erro customizado.
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// HandleError processa um erro com base na configuração.
func HandleError(err error, w http.ResponseWriter) {
	for _, config := range ErrorHandlerConfigs {
		if errors.Is(err, config.Err) {
			http.Error(w, config.ErrorCode, config.HTTPStatus)
			return
		}
	}

	// Caso não seja um erro configurado, retorna erro genérico.
	log.Printf("Unhandled error: %s", err.Error())
	http.Error(w, "internal_server_error", http.StatusInternalServerError)
}

func handleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		// Logar o erro para análise posterior
		log.Printf("Recovered from panic: %v", r)

		// Retornar uma resposta de erro genérica ao cliente
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
	}
}
