package usecase

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (u *Usecase) DeleteMusics(ids []uuid.UUID) error {
	slog.Debug("delete musics", slog.Any("ids", ids))
	if err := u.entityRepository.DeleteMusics(ids); err != nil {
		return fmt.Errorf("failed to delete musics: %w", err)
	}
	return nil
}
