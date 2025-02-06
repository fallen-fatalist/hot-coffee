package errors

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return errors.New(text)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// TODO: Add more error types
// TODO: Add more context to errors

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
