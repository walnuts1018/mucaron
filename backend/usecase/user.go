package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrIncorrectPW  = errors.New("password is incorrect")
	ErrUserExists   = errors.New("user already exists")
)

func (u *Usecase) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	return u.entityRepository.GetUserByID(ctx, id)
}

func (u *Usecase) GetUserByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.User, error) {
	return u.entityRepository.GetUserByIDs(ctx, ids)
}

func (u *Usecase) Login(ctx context.Context, userName string, inputPass entity.RawPassword) (entity.User, error) {
	user, err := u.entityRepository.GetUserByName(ctx, userName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user.LoginInfo.IsCorrectPassword(inputPass) {
		return user, nil
	} else {
		return entity.User{}, ErrIncorrectPW
	}
}

func (u *Usecase) CreateUser(ctx context.Context, userName string, inputPass entity.RawPassword) (entity.User, error) {
	validUser, err := u.IsValidUserName(ctx, userName)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to check user name: %w", err)
	}

	if !validUser {
		return entity.User{}, ErrUserExists
	}

	loginInfo, err := entity.NewLoginInfo(inputPass)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create login info: %w", err)
	}

	user, err := entity.NewUser(userName, loginInfo)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	if err := u.entityRepository.CreateUser(ctx, user); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (u *Usecase) IsValidUserName(ctx context.Context, userName string) (bool, error) {
	_, err := u.entityRepository.GetUserByName(ctx, userName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return true, nil
		}
		return false, fmt.Errorf("failed to get user by name: %w", err)
	}

	return false, nil
}
