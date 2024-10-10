package gormmodel

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDModel struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	CreatedAt synchro.Time[tz.AsiaTokyo]
	UpdatedAt synchro.Time[tz.AsiaTokyo]
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
