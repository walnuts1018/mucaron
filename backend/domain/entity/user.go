package entity

import "github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"

type User struct {
	gormmodel.UUIDModel
	UserName  string
	LoginInfo LoginInfo `gorm:"embedded"`
}

func NewUser(userName string, loginInfo LoginInfo) (User, error) {
	model, err := gormmodel.NewUUIDModel()
	if err != nil {
		return User{}, err
	}
	return User{
		UUIDModel: model,
		UserName:  userName,
		LoginInfo: loginInfo,
	}, nil
}
