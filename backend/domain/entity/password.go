package entity

import (
	"encoding/base64"

	"github.com/walnuts1018/mucaron/backend/util/random"
	"golang.org/x/crypto/scrypt"
)

type LoginInfo struct {
	HashedPassword string
	Salt           string
}

type RawPassword string

func NewLoginInfo(rawPassword RawPassword) (LoginInfo, error) {
	salt, err := random.String(32, random.AlphanumericSymbols)
	if err != nil {
		return LoginInfo{}, err
	}

	hashedPassword, err := hash(rawPassword, salt)
	if err != nil {
		return LoginInfo{}, err
	}

	return LoginInfo{
		HashedPassword: hashedPassword,
		Salt:           salt,
	}, nil
}

func hash(rawPassword RawPassword, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(rawPassword), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}

func (l LoginInfo) IsCorrectPassword(rawPassword RawPassword) bool {
	hashedPassword, err := hash(rawPassword, l.Salt)
	if err != nil {
		return false
	}
	return hashedPassword == l.HashedPassword
}
