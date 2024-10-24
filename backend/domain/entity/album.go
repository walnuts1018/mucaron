package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type Album struct {
	gormmodel.UUIDModel
	OwnerID  uuid.UUID `json:"-"`
	Owner    User      `gorm:"foreignKey:OwnerID" json:"-"`
	Name     string
	SortName *string
	Musics   []Music `gorm:"many2many:album_musics;"`
}
