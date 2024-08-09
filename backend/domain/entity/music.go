package entity

import "github.com/google/uuid"

type Music struct {
	ID               uuid.UUID   `json:"id"`
	Title            string      `json:"title"`
	AlbumID          uuid.UUID   `json:"album_id"`
	AlbumTruckNumber int64       `json:"album_truck_number"`
	ArtistIDs        []uuid.UUID `json:"artists_ids"`
	CoverInfo        CoverInfo   `json:"cover_info"`
	Score            int64       `json:"score"`
}

type CoverInfo struct {
	IsCover         bool      `json:"is_cover"`
	OriginalMusicID uuid.UUID `json:"original_music_id"`
}
