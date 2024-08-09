package entity

import "github.com/google/uuid"

type Artist struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
