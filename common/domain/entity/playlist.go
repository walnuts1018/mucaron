package entity

import (
	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/common/domain/entity/gormmodel"
)

type Playlist struct {
	gormmodel.UUIDModel
	OwnerID uuid.UUID
	Owner   User    `gorm:"foreignKey:OwnerID"`
	Musics  []Music `gorm:"many2many:playlist_musics;"`
}
