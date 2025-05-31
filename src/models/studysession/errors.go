package studysession

import "errors"

var (
	ErrActiveSessionExists   = errors.New("active session already exists")
	ErrActiveSessionNotFound = errors.New("session not found or not active")
)
