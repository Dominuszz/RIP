package repository

import (
	"context"
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	minio "lab3/internal/app/minioClient"
	"lab3/internal/app/serializer"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Repository) GetComplexClasses(searchQuery string, page, limit int) ([]ds.ComplexClass, int64, error) {
	var complexClasses []ds.ComplexClass
	var total int64

	countQuery := r.db.Model(&ds.ComplexClass{}).Where("is_delete = false")

	if searchQuery != "" {
		countQuery = countQuery.Where("degree_text ILIKE ?", "%"+searchQuery+"%")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []ds.ComplexClass{}, 0, nil
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	dataQuery := r.db.Model(&ds.ComplexClass{}).Where("is_delete = false")

	if searchQuery != "" {
		dataQuery = dataQuery.Where("degree_text ILIKE ?", "%"+searchQuery+"%")
	}

	err := dataQuery.Order("id").Offset(offset).Limit(limit).Find(&complexClasses).Error

	if err != nil {
		return nil, 0, err
	}

	return complexClasses, total, nil
}

func (r *Repository) GetComplexClass(id int) (*ds.ComplexClass, error) {
	complexClass := ds.ComplexClass{}
	err := r.db.Order("id").Where("id = ? and is_delete = ?", id, false).First(&complexClass).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w:  класс сложности  с id %d", ErrNotFound, id)
		}
		return &ds.ComplexClass{}, err
	}
	return &complexClass, nil
}

func (r *Repository) GetCompClassByDegree(degree_text string) ([]ds.ComplexClass, error) {
	items, _, err := r.GetComplexClasses(degree_text, 1, 1000)
	return items, err
}

func (r *Repository) CreateComplexClass(complexClassJSON serializer.ComplexClassJSON) (ds.ComplexClass, error) {
	complexClass := serializer.CompClassFromJSON(complexClassJSON)
	if complexClass.Degree < 0 {
		return ds.ComplexClass{}, errors.New("неправильная степень класса сложности")
	}
	err := r.db.Create(&complexClass).Error
	if err != nil {
		return ds.ComplexClass{}, err
	}
	return complexClass, nil
}

func (r *Repository) EditComplexClass(id int, complexClassJSON serializer.ComplexClassJSON) (ds.ComplexClass, error) {
	complexClass := ds.ComplexClass{}
	if id < 0 {
		return ds.ComplexClass{}, errors.New("id должно быть >= 0")
	}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&complexClass).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ComplexClass{}, fmt.Errorf("%w: класс сложности с id %d", ErrNotFound, id)
		}
		return ds.ComplexClass{}, err
	}
	if complexClassJSON.Degree < 0 {
		return ds.ComplexClass{}, errors.New("неправильная степень класса сложности")
	}
	err = r.db.Model(&complexClass).Updates(serializer.CompClassFromJSON(complexClassJSON)).Error
	if err != nil {
		return ds.ComplexClass{}, err
	}
	return complexClass, nil
}

func (r *Repository) DeleteComplexClass(id int) error {
	complexClass := ds.ComplexClass{}
	if id < 0 {
		return errors.New("id должно быть >= 0")
	}

	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&complexClass).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: класс сложности с id %d", ErrNotFound, id)
		}
		return err
	}
	if complexClass.IMG != "" {
		err = minio.DeleteObject(context.Background(), r.mc, minio.GetImgBucket(), complexClass.IMG)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.ComplexClass{}).Where("id = ?", id).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToBigORequest(bigo_request_id int, complexclass_id int) error {
	var complexClass ds.ComplexClass
	if err := r.db.First(&complexClass, complexclass_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: класс сложности с id %d", ErrNotFound, complexclass_id)
		}
		return err
	}

	var bigo_request ds.BigORequest
	if err := r.db.First(&bigo_request, bigo_request_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: заявка с id %d", ErrNotFound, bigo_request_id)
		}
		return err
	}

	comp_class_request := ds.CompClassRequest{}
	result := r.db.Where("complex_class_id = ? and big_o_request_id = ?", complexclass_id, bigo_request_id).Find(&comp_class_request)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return fmt.Errorf("%w: класс сложности %d уже в заявке %d", ErrAlreadyExists, complexclass_id, bigo_request_id)
	}
	return r.db.Create(&ds.CompClassRequest{
		ComplexClassID: uint(complexclass_id),
		BigORequestID:  uint(bigo_request_id),
	}).Error
}

func (r *Repository) GetModeratorAndCreatorLogin(bigo_request ds.BigORequest) (string, string, error) {
	var creator ds.Users
	var moderator ds.Users

	err := r.db.Where("id = ?", bigo_request.CreatorID).First(&creator).Error
	if err != nil {
		return "", "", err
	}

	var moderatorLogin string
	if bigo_request.ModeratorID.Valid {
		err = r.db.Where("id = ?", bigo_request.ModeratorID).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}

	return creator.Login, moderatorLogin, nil
}

func (r *Repository) AddPhoto(ctx *gin.Context, compclass_id int, file *multipart.FileHeader) (ds.ComplexClass, error) {
	complexclass_, err := r.GetComplexClass(compclass_id)
	if err != nil {
		return ds.ComplexClass{}, err
	}

	fileName, err := minio.UploadImage(ctx, r.mc, minio.GetImgBucket(), file, *complexclass_)
	if err != nil {
		return ds.ComplexClass{}, err
	}

	complexclass, err := r.GetComplexClass(compclass_id)
	if err != nil {
		return ds.ComplexClass{}, err
	}
	complexclass.IMG = fileName
	err = r.db.Save(&complexclass).Error
	if err != nil {
		return ds.ComplexClass{}, err
	}
	return *complexclass, nil
}
