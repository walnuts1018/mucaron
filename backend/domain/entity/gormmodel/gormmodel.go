package gormmodel

import (
	"database/sql"
	"database/sql/driver"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/uuid"
	newuuid "github.com/walnuts1018/mucaron/backend/util/new_uuid"
	"gorm.io/gorm"
)

type DeletedAt[T synchro.TimeZone] sql.Null[synchro.Time[T]]

func (n *DeletedAt[T]) Scan(value interface{}) error {
	return (*sql.Null[synchro.Time[T]])(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeletedAt[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.V, nil
}

type UUIDModel struct {
	ID        uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	CreatedAt synchro.Time[tz.AsiaTokyo]
	UpdatedAt synchro.Time[tz.AsiaTokyo]
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *UUIDModel) CreateID() error {
	id, err := newuuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}
