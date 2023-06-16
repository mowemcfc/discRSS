package models

import "errors"

var (
  ErrInternalServerError = errors.New("internal server error")
  ErrNotFound = errors.New("requested item could not be found")
	ErrConflict = errors.New("item already exists")
	ErrBadInput = errors.New("given input is not valid")
  ErrNotPermitted = errors.New("not permitted to perform given action")
)
