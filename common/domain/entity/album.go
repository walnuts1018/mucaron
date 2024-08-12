package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/common/domain/entity/gormmodel"
)

type Album struct {
	gormmodel.UUIDModel
	OwnerID  uuid.UUID
	Owner    User `gorm:"foreignKey:OwnerID"`
	Name     string
	SortName string
	Artists  []Artist `gorm:"many2many:album_artists;"`
}
