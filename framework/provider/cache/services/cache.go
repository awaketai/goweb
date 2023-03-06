package services

import (
	"errors"
	"time"
)

const NoneDuration = time.Duration(-1)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrTypeNotOk   = errors.New("val type not ok")
)