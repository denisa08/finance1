package storage

import "errors"

var (
	errURLNotFound error = errors.New("url does not found")
	errURLExists   error = errors.New("url already exists")
)
