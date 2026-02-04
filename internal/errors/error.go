package apperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Code string

const (
	CodeInternal     Code = "internal_error"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeInvalidInput Code = "invalid_input"
	CodeRateLimited  Code = "rate_limited"

	CodeChatNotFound Code = "chat_not_found"
	CodeMsgTooLong   Code = "message_too_long"
	CodeUserBlocked  Code = "user_blocked"
)

type AppError struct {
	Code    Code
	Message string            // safe for client
	Status  int               // HTTP status
	Fields  map[string]string // optional validation details
	Err     error             // internal cause
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code Code, status int, msg string) *AppError {
	return &AppError{Code: code, Status: status, Message: msg}
}

func Wrap(code Code, status int, msg string, err error) *AppError {
	return &AppError{Code: code, Status: status, Message: msg, Err: err}
}

func Is(err error, code Code) bool {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Code == code
	}
	return false
}

func WriteError(w http.ResponseWriter, err error) {
	var ae *AppError
	if !errors.As(err, &ae) {
		ae = Wrap(CodeInternal, 500, "Internal server error", err)
	}

	// log internal
	log.Printf("error code=%s status=%d err=%v", ae.Code, ae.Status, ae.Err)

	// respond public
	resp := map[string]any{
		"error": map[string]any{
			"code":    ae.Code,
			"message": ae.Message,
			"fields":  ae.Fields,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ae.Status)
	_ = json.NewEncoder(w).Encode(resp)
}
