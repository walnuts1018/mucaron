package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type Playlist struct {
	gormmodel.UUIDModel
	OwnerID uuid.UUID `json:"-"`
	Owner   User      `gorm:"foreignKey:OwnerID" json:"-"`
	Musics  []Music   `gorm:"many2many:playlist_musics;"`
}
