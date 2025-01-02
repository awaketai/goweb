package config

import "errors"

var (
	KeyNotFoundError      = errors.New("the key is not found")
	InvalidValueTypeError = errors.New("the value is not expected type")
)
