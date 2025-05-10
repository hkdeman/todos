package todos

import "errors"

var (
	ErrTodoNotFound     = errors.New("todo not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidDate      = errors.New("invalid date")
	ErrInvalidPriority  = errors.New("invalid priority")
)
