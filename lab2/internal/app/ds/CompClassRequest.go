package ds

type CompClassRequest struct {
	ID             uint    `gorm:"primaryKey"`
	BigORequestID  uint    `gorm:"not null;uniqueIndex:idx_compclass_request"`
	ComplexClassID uint    `gorm:"not null;uniqueIndex:idx_compclass_request"`
	Complexity     string  `gorm:"type:varchar(100);not null"`
	Degree         float64 `gorm:"not null"`
	DegreeText     string  `gorm:"type:varchar(100); not null"`
	ArraySize      uint
	BigORequest    BigORequest  `gorm:"foreignKey:BigORequestID;references:ID"`
	ComplexClass   ComplexClass `gorm:"foreignKey:ComplexClassID;references:ID"`
}
