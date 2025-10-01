package ds

type ComplexClass struct {
	ID          uint   `gorm:"primaryKey"`
	IMG         string `gorm:"type:varchar(100)"`
	Complexity  string `gorm:"type:varchar(100);not null"`
	Degree      string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:varchar(255)"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false"`
}
