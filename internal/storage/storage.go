package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

const (
	UniqueViolationCode = "23505"
)
