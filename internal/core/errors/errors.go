package errors

import (
	"errors"
	"fmt"
)

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

func New(text string) error {
	return errors.New(text)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
