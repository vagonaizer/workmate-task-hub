package models

import "errors"

var (
	ErrInvalidTitle    = errors.New("title is required")
	ErrInvalidStatus   = errors.New("invalid status")
	ErrInvalidDeadline = errors.New("invalid deadline")
	ErrInvalidPriority = errors.New("invalid priority")
)
