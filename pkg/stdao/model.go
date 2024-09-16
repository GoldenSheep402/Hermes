package stdao

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        string `gorm:"primary_key;type:char(26)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Model) BeforeCreate(_ *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}
	return
}
