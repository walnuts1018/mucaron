package gormmodel

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt synchro.Time[tz.AsiaTokyo]
	UpdatedAt synchro.Time[tz.AsiaTokyo]
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UUIDModel struct {
	BaseModel
	ID uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
}

type HashModel struct {
	BaseModel
	ID string `gorm:"primarykey"`
}
