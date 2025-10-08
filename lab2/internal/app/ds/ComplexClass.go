package ds

type ComplexClass struct {
	ID          uint    `gorm:"primaryKey"`
	IMG         string  `gorm:"type:varchar(100)"`
	Complexity  string  `gorm:"type:varchar(100);not null"`
	Degree      float64 `gorm:"not null"`
	DegreeText  string  `gorm:"type:varchar(100); not null"`
	Description string  `gorm:"type:varchar(255); not null"`
	IsDelete    bool    `gorm:"type:boolean not null;default:false"`
}
