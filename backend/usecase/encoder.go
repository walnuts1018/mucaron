package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
)

func (u *Usecase) encodeLockedSection(ctx context.Context, uploadedFilePath string, music entity.Music) (string, error) {
	u.encodeMutex.Lock()
	defer u.encodeMutex.Unlock()

	if err := u.entityRepository.UpdateMusicStatus(music.ID, entity.VideoEncoding); err != nil {
		return "", fmt.Errorf("failed to update music status: %w", err)
	}

	slog.Info("start encoding", slog.String("music_id", music.ID.String()))

	outDir, err := u.encoder.Encode(music.ID.String(), uploadedFilePath, false)
	if err != nil {
		return "", fmt.Errorf("failed to encode: %w", err)
	}
	slog.Info("finish encoding", slog.String("music_id", music.ID.String()))
	return outDir, nil
}

func (u *Usecase) encode(ctx context.Context, uploadedFilePath string, music entity.Music) {
	outDir, err := u.encodeLockedSection(ctx, uploadedFilePath, music)
	if err != nil {
		slog.Error("failed to encode", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		// status更新
		return
	}

	encodedUploadedFilePath := strings.ReplaceAll(uploadedFilePath, uploadedExtension, encodedExtension)
	if err := os.Rename(uploadedFilePath, encodedUploadedFilePath); err != nil {
		slog.Error("failed to rename uploaded file", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		// ファイル移動に失敗しても、エンコードは成功しているので、returnしない
	}

	primaryM3U8, err := os.ReadFile(filepath.Join(outDir, "primary.m3u8"))
	if err != nil {
		slog.Error("failed to open primary.m3u8", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}

	replaceM3U8URL, err := replaceM3U8URL(string(primaryM3U8), u.cfg.ServerEndpoint, music.ID.String())
	if err != nil {
		slog.Error("failed to replace m3u8 url", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}
	// os.WriteFileにはO_TRUNCフラグがあるので、すでにファイルが存在する場合、上書きされる
	if err := os.WriteFile(filepath.Join(outDir, "primary.m3u8"), []byte(replaceM3U8URL), 0644); err != nil {
		slog.Error("failed to write primary.m3u8", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}

	slog.Info("start uploading", slog.String("music_id", music.ID.String()))
	if err := u.objectStorage.UploadDirectory(ctx, music.ID.String(), outDir); err != nil {
		slog.Error("failed to upload directory", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		return
	}
	slog.Info("finish uploading", slog.String("music_id", music.ID.String()))

	if err := os.RemoveAll(outDir); err != nil {
		slog.Error("failed to remove directory", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		// ファイル削除に失敗しても、アップロードは成功しているので、returnしない
	}

	if err := os.Remove(encodedUploadedFilePath); err != nil {
		slog.Error("failed to remove uploaded file", slog.String("music_id", music.ID.String()), slog.Any("error", err))
		// ファイル削除に失敗しても、アップロードは成功しているので、returnしない
	}

	if err := u.entityRepository.UpdateMusicStatus(music.ID, entity.VideoEncoded); err != nil {
		slog.Error("failed to update music status", slog.Any("error", err), slog.String("music_id", music.ID.String()))
	}
}

func (u *Usecase) EncodeSuspended(ctx context.Context) error {
	files, err := os.ReadDir(os.TempDir())
	if err != nil {
		return err
	}

	musics := make([]entity.Music, 0, len(files))
	filePathes := make([]string, 0, len(files))
	for _, file := range files {
		slog.Debug("found candidate file", slog.String("file_name", file.Name()))

		if file.IsDir() {
			if strings.Contains(file.Name(), u.encoder.GetOutDirPrefix()) {
				slog.Debug("found outdir", slog.String("file_name", file.Name()))
				if err := os.RemoveAll(filepath.Join(os.TempDir(), file.Name())); err != nil {
					slog.Error("failed to remove directory", slog.String("file_name", file.Name()), slog.Any("error", err))
				}
			}
			slog.Debug("skip directory", slog.String("file_name", file.Name()))
			continue
		}

		if !strings.Contains(filepath.Ext(file.Name()), uploadedExtension) {
			slog.Debug("skip non-uploaded file", slog.String("file_name", file.Name()))
			continue
		}

		id, err := uuid.Parse(file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))])
		if err != nil {
			slog.Error("failed to parse uuid", slog.String("file_name", file.Name()), slog.Any("error", err))
			// 無効なファイルを削除
			if err := os.Remove(filepath.Join(os.TempDir(), file.Name())); err != nil {
				slog.Error("failed to remove invalid file", slog.String("file_name", file.Name()), slog.Any("error", err))
			}
			continue
		}

		music, err := u.entityRepository.GetMusicByID(id)
		if err != nil {
			slog.Error("failed to get music by id", slog.String("music_id", id.String()), slog.Any("error", err))
			// すでに削除されている / 無効なファイルを削除
			// このメソッドはサーバー起動前にのみ実行されることを想定しているので、ファイルがアップロードされたもののまだDBには存在していない、というケースは無視する
			if err := os.Remove(filepath.Join(os.TempDir(), file.Name())); err != nil {
				slog.Error("failed to remove invalid file", slog.String("file_name", file.Name()), slog.Any("error", err))
			}
			continue
		}

		if music.Status != entity.VideoEncoded {
			slog.Debug("skip encoded music", slog.String("music_id", music.ID.String()))
			// すでにエンコード済みのファイルを削除
			if err := os.Remove(filepath.Join(os.TempDir(), file.Name())); err != nil {
				slog.Error("failed to remove invalid file", slog.String("file_name", file.Name()), slog.Any("error", err))
			}
			continue
		}

		musics = append(musics, music)
		filePathes = append(filePathes, filepath.Join(os.TempDir(), file.Name()))
	}

	if len(musics) == 0 {
		return nil
	}

	musicIDs := make([]uuid.UUID, len(musics))
	for i, music := range musics {
		musicIDs[i] = music.ID
	}
	if err := u.entityRepository.UpdateMusicStatuses(musicIDs, entity.EncodeRetrying); err != nil {
		return err
	}

	for i, music := range musics {
		filePath := filePathes[i]
		slog.Info("start encoding suspended music", slog.String("music_id", music.ID.String()), slog.String("file_path", filePath))
		go func() {
			u.encode(ctx, filePath, music)
		}()
	}
	return nil
}
