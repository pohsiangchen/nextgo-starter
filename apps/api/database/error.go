package database

import "errors"

var (
	ErrDuplicateEmail = errors.New("a user with that email already exists")
	ErrRecordNotFound = errors.New("record not found")
)
