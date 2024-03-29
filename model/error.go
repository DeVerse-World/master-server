package model

import (
	"errors"
)

// Model errors
var (
	ErrInvalidArgs  = errors.New("Invalid Args")
	ErrKeyConflict  = errors.New("Key Conflict")
	ErrDataNotFound = errors.New("Record Not Found")
	ErrUserExists   = errors.New("User already exists")
	ErrUnknown      = errors.New("Unknown Error")
	ErrFailed       = errors.New("Failed")
	ErrExpiredToken = errors.New("Your access token expired")
	ErrInvalidToken = errors.New("JWT Token invalid")
)

var (
	DbDuplicateEntryCode = uint16(1062)
)
