package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type Genre struct {
	gormmodel.UUIDModel
	OwnerID uuid.UUID `json:"-"`
	Owner   User      `gorm:"foreignKey:OwnerID" json:"-"`
	Name    string    `json:"name"`
}
