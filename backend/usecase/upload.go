package usecase

import (
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/util/result"
)

func (u *Usecase) UploadMusic(user entity.User, r io.Reader) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}
	ch := make(chan result.Result[string])

	go func(ch chan<- result.Result[string]) {
		path, err := u.encoder.Encode(id, r, false)
		ch <- result.Result[string]{
			Result: path,
			Error:  err,
		}
	}(ch)

	m := entity.Music{
		OwnerID: user.ID,
	}

	if err := u.MusicRepository.CreateMusic(m); err != nil {
		return fmt.Errorf("failed to create music: %w", err)
	}
}
