package domain

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = gorm.ErrRecordNotFound
var ErrAlreadyExists = errors.New("music already exists")