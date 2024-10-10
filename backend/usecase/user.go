package usecase

import (
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (u *Usecase) GetUserByID(id string) (entity.User, error) {
	return u.userRepository.GetUserByID(id)
}

func (u *Usecase) GetUserByIDs(ids []string) ([]entity.User, error) {
	return u.userRepository.GetUserByIDs(ids)
}

func (u *Usecase) Login(userID uuid.UUID, pass entity.RawPassword) (entity.User, error) {
	user, err := u.userRepository.GetUserByID(userID)
	if err != nil {
		return entity.User{}, err
	}

