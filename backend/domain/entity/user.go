package entity

import "github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"

type User struct {
	gormmodel.UUIDModel
	UserName  string
	LoginInfo LoginInfo `gorm:"embedded" json:"-"`
}

