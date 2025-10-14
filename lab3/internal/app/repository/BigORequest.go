package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"
	"math"
	"time"

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
	sub := r.db.Table("comp_class_requests").Where("big_o_request_id = ?", bigorequest.ID)
	err = r.db.Order("id DESC").Where("id IN (?)", sub.Select("id")).Find(&compclasses).Error

	if err != nil {
		return []ds.ComplexClass{}, ds.BigORequest{}, err
	}

	return compclasses, bigorequest, nil
}

func (r *Repository) CheckCurrentBigORequestDraft(creator_ID uint) (ds.BigORequest, error) {
	var bigorequest ds.BigORequest
	res := r.db.Where("creator_id = ? AND status = ?", creator_ID, "черновик").Limit(1).Find(&bigorequest)
	if res.Error != nil {
		return ds.BigORequest{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.BigORequest{}, ErrNoDraft
	}
	return bigorequest, nil
}

func (r *Repository) GetBigORequestDraft(creator_ID uint) (ds.BigORequest, bool, error) {
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

func (r *Repository) GetBigORequestCount(creator_ID uint) int64 {
	if creator_ID == 0 {
		return 0
	}

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

func (r *Repository) FinishBigORequest(id int, status string) (ds.BigORequest, error) {
	if status != "выполнен" && status != "отклонен" {
		return ds.BigORequest{}, errors.New("неверный статус")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.BigORequest{}, err
	}

	if !user.IsModerator {
		return ds.BigORequest{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}

	bigorequest, err := r.GetSingleBigORequest(id)
	if err != nil {
		return ds.BigORequest{}, err
	} else if bigorequest.Status != "сформирован" {
		return ds.BigORequest{}, fmt.Errorf("это исследование не может быть %s", status)
	}

	err = r.db.Model(&bigorequest).Updates(ds.BigORequest{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: uint(user.ID),
	}).Error
	if err != nil {
		return ds.BigORequest{}, err
	}
	if status == "выполнен" {
		var res = 0.0
		compclassrequest, err := r.GetComplexClassesBigORequests(int(bigorequest.ID))
		if err != nil {
			return ds.BigORequest{}, err
		}
		for _, compclassrequest := range compclassrequest {
			compclass_time, err := CalculateComplexClassTime(compclassrequest.Degree, compclassrequest.ArraySize)
			if err != nil {
				return ds.BigORequest{}, err
			}
			res += compclass_time
		}
		bigorequest.CalculatedTime = res
	}

	return bigorequest, nil
}
