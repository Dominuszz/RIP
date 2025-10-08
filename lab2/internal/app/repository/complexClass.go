package repository

import (
	"errors"
	"fmt"
	"time"

	"lab2/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetComplexClasses() ([]ds.ComplexClass, error) {
	var ComplexClasses []ds.ComplexClass
	err := r.db.Find(&ComplexClasses).Error
	if err != nil {
		return nil, err
	}
	if len(ComplexClasses) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return ComplexClasses, nil
}

func (r *Repository) GetComplexClass(id int) (ds.ComplexClass, error) {
	complexclass := ds.ComplexClass{}
	err := r.db.Where("ID = ?", id).Find(&complexclass).Error
	if err != nil {
		return ds.ComplexClass{}, err
	}
	return complexclass, nil
}

func (r *Repository) GetComplexClasssByDegree(title string) ([]ds.ComplexClass, error) {
	var ComplexClasses []ds.ComplexClass
	err := r.db.Where("degree_text ILIKE ?", "%"+title+"%").Find(&ComplexClasses).Error
	if err != nil {
		return nil, err
	}
	return ComplexClasses, nil
}
func (r *Repository) GetActiveBigORequestID() uint {
	var BigORequestID uint
	err := r.db.Model(&ds.BigORequest{}).Where("status = ?", "черновик").Select("id").First(&BigORequestID).Error
	if err != nil {
		return 0
	}
	return BigORequestID
}

func (r *Repository) GetBigORequestCount() int64 {
	var bigorequestID uint
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.BigORequest{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("id").First(&bigorequestID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CompClassRequest{}).Where("big_o_request_id = ?", bigorequestID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_chats:", err)
	}

	return count
}

func (r *Repository) AddComplexClass(ComplexClassID uint, creatorID uint) error {
	var request ds.BigORequest

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").
		First(&request).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		request = ds.BigORequest{
			Status:      "черновик",
			DateCreate:  time.Now(),
			CreatorID:   creatorID,
			ModeratorID: 2,
		}
		if err := r.db.Create(&request).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.CompClassRequest{}).
		Where("big_o_request_id = ? AND complex_class_id = ?", request.ID, ComplexClassID).Preload("ComplexClass").
		Count(&count)

	if count == 0 {
		var complexclass ds.ComplexClass
		if err := r.db.First(&complexclass, ComplexClassID).Error; err != nil {
			return err
		}

		appDev := ds.CompClassRequest{
			BigORequestID:  request.ID,
			ComplexClassID: ComplexClassID,
			Complexity:     complexclass.Complexity,
			Degree:         complexclass.Degree,
		}
		if err := r.db.Create(&appDev).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteBigORequest(RequestID uint) error {
	query := `
		UPDATE big_o_requests
		SET status = 'удалён'
		WHERE id = $1;
	`
	result := r.db.Exec(query, RequestID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application with id %d not found", RequestID)
	}
	return nil
}

func (r *Repository) IsDraftBigORequest(RequestID int) (bool, error) {
	var request ds.BigORequest
	err := r.db.Select("status").Where("id = ?", RequestID).First(&request).Error
	if err != nil {
		return false, err
	}
	return request.Status == "черновик", nil
}

func (r *Repository) GetBigORequest(id int) ([]ds.CompClassRequest, error) {
	var RequestItems []ds.CompClassRequest
	err := r.db.Where("big_o_request_id= ?", id).Preload("ComplexClass").Find(&RequestItems).Error
	if err != nil {
		return nil, err
	}

	return RequestItems, nil
}
