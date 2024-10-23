package newuuid

import "github.com/google/uuid"

var testValue uuid.UUID = uuid.Nil

func SetUUIDValue(value uuid.UUID) {
	testValue = value
}

func NewV7() (uuid.UUID, error) {
	if testValue != uuid.Nil {
		return testValue, nil
	} else {
		return uuid.NewV7()
	}
}
