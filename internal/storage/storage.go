package storage

import "errors"

var(
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExists = errors.New("url exists")
	ErrUrlNotDeleted = errors.New("no rows deleted")
)

