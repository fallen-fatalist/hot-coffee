package serviceinstance

import (
	"errors"
	"strconv"
)

// errors
var (
	// ID errors
	ErrEmptyID      = errors.New("empty id provided")
	ErrNonNumericID = errors.New("non-numeric id provided")
	ErrNegativeID   = errors.New("negative id provided")
	ErrZeroID       = errors.New("zero id provided")
)

// Function checks ID to be positive integer
func isValidID(id string) error {
	if id == "" {
		return ErrEmptyID
	}

	itemID, err := strconv.Atoi(id)
	if err != nil {
		return ErrNonNumericID
	} else if itemID == 0 {
		return ErrZeroID
	} else if itemID < 1 {
		return ErrNegativeID
	}

	return nil
}
