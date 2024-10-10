package entity

type HashedPassword string
type RawPassword string

type HashSalt string

func NewHashSalt() HashSalt {
	return HashSalt("dummy")
}
