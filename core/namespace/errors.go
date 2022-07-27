package namespace

import "errors"

var (
	ErrInvalidID = errors.New("namespace id is invalid")
	ErrNotExist  = errors.New("namespace doesn't exist")
	ErrConflict  = errors.New("namespace name already exist")
)
