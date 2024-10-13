package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type Artist struct {
	gormmodel.UUIDModel
	OwnerID  uuid.UUID
	Owner    User    `gorm:"foreignKey:OwnerID"`
	Name     string  `json:"name"`
	SortName *string `json:"sort_name"`
}
