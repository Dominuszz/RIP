package ds

import "github.com/google/uuid"

type Users struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Login       string    `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string    `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool      `gorm:"type:boolean;default:false" json:"is_moderator"`
}
