package ds

import (
	"database/sql"
	"time"
)

type BigORequest struct {
	ID                   uint      `gorm:"primaryKey"`
	Status               string    `gorm:"type:varchar(15);not null"`
	DateCreate           time.Time `gorm:"not null"`
	DateUpdate           time.Time
	DateFinish           sql.NullTime `gorm:"default:null"`
	CreatorID            uint         `gorm:"not null"`
	ModeratorID          uint
	CalculatedTime       float64 `gorm:"type:numeric(3,1)"`
	CalculatedComplexity string  `gorm:"type:varchar(255)"`
	Creator              Users   `gorm:"foreignKey:CreatorID"`
	Moderator            Users   `gorm:"foreignKey:ModeratorID"`
}
