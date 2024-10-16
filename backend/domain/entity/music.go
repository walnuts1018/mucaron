package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

type Music struct {
	gormmodel.UUIDModel
	OwnerID          uuid.UUID `gorm:"index:idx_musics_file_hash,unique" json:"-"`
	Owner            User      `gorm:"foreignKey:OwnerID;" json:"-"`
	Name             string
	SortName         *string
	AlbumTrackNumber *int64
	Artists          []Artist `gorm:"many2many:music_artists;"`
	Score            int64
	Duration         time.Duration
	GenreID          *uuid.UUID
	Genre            *Genre

	// ----- cover info -----
	IsCover         bool
	OriginalMusicID *uuid.UUID
	// ----------------------

	// ----- raw data -----
	RawMetaData RawMusicMetadata `json:"-"`

	// ----- status -----
	Status    MusicStatus
	PlayCount int64

	FileHash string `gorm:"index:idx_musics_file_hash,unique"`
}

type MusicStatus string

const (
	Unknown             MusicStatus = "unknown"
	MetadataParsed      MusicStatus = "metadata_parsed"
	EncodeRetrying      MusicStatus = "encode_retrying"
	AudioEncoding       MusicStatus = "audio_encoding"
	AudioEncodeFinished MusicStatus = "audio_encode_finished"
	VideoEncoding       MusicStatus = "video_encoding"
	VideoEncoded        MusicStatus = "video_encoded"
)
