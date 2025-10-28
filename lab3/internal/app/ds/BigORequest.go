package ds

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BigORequest struct {
	ID                   uint          `gorm:"primary_key;autoIncrement"`
	Status               string        `gorm:"type:varchar(15);not null"`
	DateCreate           time.Time     `gorm:"not null"`
	DateUpdate           sql.NullTime  `gorm:"default:null"`
	DateFinish           sql.NullTime  `gorm:"default:null"`
	CreatorID            uuid.UUID     `gorm:"not null"`
	ModeratorID          uuid.NullUUID `gorm:"default:null"`
	CalculatedTime       float64       `gorm:"type:numeric(30,10)"`
	CalculatedComplexity string        `gorm:"type:varchar(255)"`
	Creator              Users         `gorm:"foreignKey:CreatorID"`
	Moderator            Users         `gorm:"foreignKey:ModeratorID"`
}
