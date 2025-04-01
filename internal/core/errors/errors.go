package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// TODO: Add more error types
// TODO: Add more context to errors

// Errors instances
var (
	ErrIncorrectRequest = New("incorrect request provided")
	ErrIDAlreadyExists  = New("entity with such id already exists")
)

// General Application error type \\
type appError struct {
	Error   error
	Message string
	Code    int
}

// Error enveloper for Encoding to JSON \\
type errorEnveloper struct {
	Message string `json:"error"`
}

// Handler Returning some App error \\
type appHandler func(http.ResponseWriter, *http.Request) *appError

// Middleware for handling error \\
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if appErr := fn(w, r); appErr != nil {
		slog.Error(appErr.Error.Error())
		json, _ := json.Marshal(errorEnveloper{Message: appErr.Message})

		w.WriteHeader(appErr.Code)
		w.Write(json)

	}
}

// MUST DO: Add middleware for catching the panics

type ErrInsufficientIngredient struct {
	insufficientIngredientID string
}

func NewErrInsufficientIngredient(ingridientID string) *ErrInsufficientIngredient {
	return &ErrInsufficientIngredient{
		ingridientID,
	}
}

func (err *ErrInsufficientIngredient) Error() string {
	return fmt.Sprintf("Insufficient: %s ingredient id", err.insufficientIngredientID)
}

func (err *ErrInsufficientIngredient) Unwrap() error {
	return nil
}

type ErrNonIntegerID struct {
	entity     string
	providedID string
}

func NewErrNonIntegerID(entity, providedID string) *ErrNonIntegerID {
	return &ErrNonIntegerID{entity, providedID}
}

func (err *ErrNonIntegerID) Error() string {
	return fmt.Sprintf("%s id of %s is not integer", err.providedID, err.entity)
}

// Wrappers over standard errors package \\

func New(text string) error {
	return errors.New(text)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
