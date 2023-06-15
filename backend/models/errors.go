package models

import "errors"

var (
  ErrInternalServerError = errors.New("Internal Server Error")
  ErrNotFound = errors.New("requested item could not be found")
	ErrConflict = errors.New("item already exists")
	ErrBadParamInput = errors.New("given param is not valid")
  ErrNotPermitted = errors.New("not permitted to perform given action")
)
