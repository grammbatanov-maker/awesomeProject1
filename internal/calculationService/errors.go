package calculationService

import "errors"

var (
	ErrInvalidID           = errors.New("invalid calculation id")
	ErrCalculationNotFound = errors.New("calculation not found")
	ErrInvalidExpression   = errors.New("expression cannot be empty")
)
