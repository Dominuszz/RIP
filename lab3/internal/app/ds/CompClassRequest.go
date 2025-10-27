package ds

type CompClassRequest struct {
	ID             uint `gorm:"primary_key;autoIncrement"`
	BigORequestID  uint `gorm:"not null;uniqueIndex:idx_compclass_request"`
	ComplexClassID uint `gorm:"not null;uniqueIndex:idx_compclass_request"`
	ArraySize      uint
	BigORequest    BigORequest  `gorm:"foreignKey:BigORequestID;references:ID"`
	ComplexClass   ComplexClass `gorm:"foreignKey:ComplexClassID;references:ID"`
}
