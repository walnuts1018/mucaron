package entity

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/walnuts1018/mucaron/backend/util/random"
	"golang.org/x/crypto/scrypt"
)

type LoginInfo struct {
	HashedPassword string
	Salt           string
}

func NewLoginInfo(rawPassword RawPassword) (LoginInfo, error) {
	if err := IsValidPassword(rawPassword); err != nil {
		return LoginInfo{}, err
	}

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

type RawPassword string

var (
	ErrInvalidPassword = errors.New("invalid password")

	ErrPasswordLength = fmt.Errorf("password length must be 8 to 128: %w", ErrInvalidPassword)
	ErrPasswordFormat = fmt.Errorf("password must include at least 2 of the following: lowercase letters, uppercase letters, numbers, and symbols: %w", ErrInvalidPassword)
)

func IsValidPassword(rawPassword RawPassword) error {
	if !(8 <= len(rawPassword) && len(rawPassword) <= 128) {
		return ErrPasswordLength
	}

	var hasLowerAlpha, hasUpperAlpha, hasNumber, hasOther int
	for _, r := range rawPassword {
		switch getCharType(r) {
		case CharTypeLowerAlpha:
			hasLowerAlpha = 1
		case CharTypeUpperAlpha:
			hasUpperAlpha = 1
		case CharTypeNumber:
			hasNumber = 1
		case CharTypeOther:
			hasOther = 1
		}
	}

	if hasLowerAlpha+hasUpperAlpha+hasNumber+hasOther < 2 {
		return ErrPasswordFormat
	}

	return nil
}

type CharType string

const (
	//アルファベット小文字
	CharTypeLowerAlpha CharType = "lower_alpha"
	//アルファベット大文字
	CharTypeUpperAlpha CharType = "upper_alpha"
	//数字
	CharTypeNumber CharType = "number"
	//その他文字
	CharTypeOther CharType = "other"
)

func getCharType(r rune) CharType {
	switch {
	case 'a' <= r && r <= 'z':
		return CharTypeLowerAlpha
	case 'A' <= r && r <= 'Z':
		return CharTypeUpperAlpha
	case '0' <= r && r <= '9':
		return CharTypeNumber
	default:
		return CharTypeOther
	}
}
