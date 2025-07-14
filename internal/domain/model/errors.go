package model

import "errors"

// Domain errors - simple and focused
var (
	// Pack validation errors
	ErrInvalidPackSize   = errors.New("pack size must be greater than zero")
	ErrInvalidPackName   = errors.New("pack name cannot be empty")
	ErrPackNotFound      = errors.New("pack not found")
	ErrPackAlreadyExists = errors.New("pack already exists")

	// Calculation errors
	ErrInvalidOrderQuantity = errors.New("order quantity must be greater than zero")
	ErrEmptyPackSizes       = errors.New("pack sizes cannot be empty")
	ErrCalculationFailed    = errors.New("unable to calculate pack distribution")

	// Business rule errors
	ErrNoValidPacks  = errors.New("no valid pack configurations available")
	ErrOrderTooLarge = errors.New("order quantity exceeds maximum limit")
)
