package storage

import "errors"

var (
	ErrURLNotFound     = errors.New("url not found")
	ErrShortCodeExists = errors.New("short code already exists")
)
