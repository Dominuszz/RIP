package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var errNoDraft = errors.New("no draft for this user")

func (r *Repository) GetAllBigORequests(from, to time.Time, status string) ([]ds.BigORequest, error) {
	var bigo_requests []ds.BigORequest
	sub := r.db.Where("status != 'удален' and status != 'черновик'")
	if !from.IsZero() {
		sub = sub.Where("date_create > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_create < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("id").Find(&bigo_requests).Error
	if err != nil {
		return nil, err
	}
	return bigo_requests, nil
}

func (r *Repository) GetComplexClassesBigORequests(bigo_request_id int) ([]ds.CompClassRequest, error) {
	var compclassrequest []ds.CompClassRequest
	err := r.db.Where("big_o_request_id = ?", bigo_request_id).Find(&compclassrequest).Error
	if err != nil {
		return nil, err
	}
	return compclassrequest, nil
}

func (r *Repository) GetComplexClassesBigORequest(compclass_id int, bigo_request_id int) (ds.CompClassRequest, error) {
	var compclassrequest ds.CompClassRequest
	err := r.db.Where("complex_class_id = ? and big_o_request_id = ?", compclass_id, bigo_request_id).First(&compclassrequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.CompClassRequest{}, fmt.Errorf("%w: complex_class bigo_request not found", ErrNotFound)
		}
		return ds.CompClassRequest{}, err
	}
	return compclassrequest, nil
}

func (r *Repository) GetBigORequestComplexClasses(id int) ([]ds.ComplexClass, ds.BigORequest, error) {
	bigorequest, err := r.GetSingleBigORequest(id)
	if err != nil {
		return []ds.ComplexClass{}, ds.BigORequest{}, err
	}

	var compclasses []ds.ComplexClass
	err = r.db.
		Joins("JOIN comp_class_requests ON comp_class_requests.complex_class_id = complex_classes.id").
		Where("comp_class_requests.big_o_request_id = ?", bigorequest.ID).
		Order("complex_classes.id DESC").
		Find(&compclasses).Error

	if err != nil {
		return []ds.ComplexClass{}, ds.BigORequest{}, err
	}

	return compclasses, bigorequest, nil
}

func (r *Repository) CheckCurrentBigORequestDraft(creator_ID uuid.UUID) (ds.BigORequest, error) {
	var bigorequest ds.BigORequest
	res := r.db.Where("creator_id = ? AND status = ?", creator_ID, "черновик").Limit(1).Find(&bigorequest)
	if res.Error != nil {
		return ds.BigORequest{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.BigORequest{}, ErrNoDraft
	}
	return bigorequest, nil
}

func (r *Repository) GetBigORequestDraft(creator_ID uuid.UUID) (ds.BigORequest, bool, error) {
	bigorequest, err := r.CheckCurrentBigORequestDraft(creator_ID)
	if errors.Is(err, ErrNoDraft) {
		bigorequest = ds.BigORequest{
			Status:     "черновик",
			CreatorID:  creator_ID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&bigorequest)
		if result.Error != nil {
			return ds.BigORequest{}, false, result.Error
		}
		return bigorequest, true, nil
	} else if err != nil {
		return ds.BigORequest{}, false, err
	}
	return bigorequest, true, nil
}

func (r *Repository) GetBigORequestCount(creator_ID uuid.UUID) int64 {
	var count int64
	bigorequest, err := r.CheckCurrentBigORequestDraft(creator_ID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.CompClassRequest{}).Where("big_o_request_id = ?", bigorequest.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_devices:", err)
	}

	return count
}

func (r *Repository) DeleteCalculation(bigo_request_id int) error {
	return r.db.Exec("UPDATE big_o_requests SET status = 'удален' WHERE id = ?", bigo_request_id).Error
}

func (r *Repository) GetSingleBigORequest(id int) (ds.BigORequest, error) {
	if id < 0 {
		return ds.BigORequest{}, errors.New("неверное id, должно быть >= 0")
	}

	var bigo_request ds.BigORequest
	err := r.db.Where("id = ?", id).First(&bigo_request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.BigORequest{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.BigORequest{}, err
	} else if bigo_request.Status == "удален" {
		return ds.BigORequest{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return bigo_request, nil
}

func (r *Repository) FormBigORequest(bigo_request_id int, status string) (ds.BigORequest, error) {
	bigorequest, err := r.GetSingleBigORequest(bigo_request_id)
	if err != nil {
		return ds.BigORequest{}, err
	}

	if bigorequest.Status != "черновик" {
		return ds.BigORequest{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	err = r.db.Model(&bigorequest).Updates(ds.BigORequest{
		Status: status,
		DateUpdate: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.BigORequest{}, err
	}

	return bigorequest, nil
}

func (r *Repository) EditBigORequest(id int, bigo_requestJSON serializer.BigORequestJSON) (ds.BigORequest, error) {
	bigorequest := ds.BigORequest{}
	if id < 0 {
		return ds.BigORequest{}, errors.New("неправильное id, должно быть >= 0")
	}
	if bigorequest.CalculatedTime < 0 {
		return ds.BigORequest{}, errors.New("неправильная нагрузка")
	}
	err := r.db.Where("id = ? and status != 'удален'", id).First(&bigorequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.BigORequest{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.BigORequest{}, err
	}
	err = r.db.Model(&bigorequest).Updates(serializer.BigORequestFromJSON(bigo_requestJSON)).Error
	if err != nil {
		return ds.BigORequest{}, err
	}
	return bigorequest, nil
}

func CalculateComplexClassTime(degree float64, arraysize uint) (float64, error) {
	if arraysize < 0 {
		return 0, errors.New("неправильная длина массива")
	}
	if degree < 0 {
		return 0, errors.New("неправильная степень класса сложности")
	}
	return math.Pow(cast.ToFloat64(arraysize), degree), nil
}

func (r *Repository) FinishBigORequest(id int, status string, currentUserID uuid.UUID) (ds.BigORequest, error) {
	if status != "выполнен" && status != "отклонен" {
		return ds.BigORequest{}, errors.New("неверный статус")
	}

	bigorequest, err := r.GetSingleBigORequest(id)
	if err != nil {
		return ds.BigORequest{}, err
	} else if bigorequest.Status != "сформирован" {
		return ds.BigORequest{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	// Только меняем статус и модератора, БЕЗ РАСЧЕТА
	err = r.db.Model(&bigorequest).Updates(map[string]interface{}{
		"status":       status,
		"date_finish":  time.Now(),
		"moderator_id": currentUserID,
		"date_update":  time.Now(),
	}).Error

	if err != nil {
		return ds.BigORequest{}, err
	}

	return bigorequest, nil
}

// UpdateBigORequestResult обновляет результаты расчета заявки
func (r *Repository) UpdateBigORequestResult(bigorequest ds.BigORequest) error {
	updates := map[string]interface{}{
		"calculated_complexity": bigorequest.CalculatedComplexity,
		"calculated_time":       bigorequest.CalculatedTime,
		"date_update":           time.Now(),
	}

	// Обновляем статус, если он изменился
	if bigorequest.Status != "" {
		updates["status"] = bigorequest.Status
	}

	return r.db.Model(&bigorequest).Updates(updates).Error
}
