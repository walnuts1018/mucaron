package usecase

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (u *Usecase) DeleteMusics(user entity.User, ids []uuid.UUID) error {
	slog.Debug("delete musics", slog.Any("ids", ids))
	if err := u.entityRepository.DeleteMusics(ids); err != nil {
		return fmt.Errorf("failed to delete musics: %w", err)
	}
	return nil
}

func (u *Usecase) GetMusics(user entity.User) ([]entity.Music, error) {
	musics, err := u.entityRepository.GetMusicsByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get musics: %w", err)
	}
	return musics, nil
}

func (u *Usecase) GetMusicIDs(user entity.User) ([]uuid.UUID, error) {
	ids, err := u.entityRepository.GetMusicIDsByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get music ids: %w", err)
	}
	return ids, nil
}
