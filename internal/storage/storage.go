package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
