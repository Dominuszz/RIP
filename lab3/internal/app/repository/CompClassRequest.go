package repository

import (
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"

	"gorm.io/gorm"
)

func (r *Repository) DeleteCompClassFromBigORequest(bigo_request_id int, compclass_id int) (ds.BigORequest, error) {
	var bigo_request ds.BigORequest
	err := r.db.Where("id = ?", bigo_request_id).First(&bigo_request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.BigORequest{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, bigo_request_id)
		}
		return ds.BigORequest{}, err
	}
	err = r.db.Where("complex_class_id = ? and big_o_request_id = ?", compclass_id, bigo_request_id).Delete(&ds.CompClassRequest{}).Error
	if err != nil {
		return ds.BigORequest{}, err
	}
	return bigo_request, nil
}

func (r *Repository) EditCompClassFromBigORequest(big_o_request_id int, compclass_id int, CompClassRequestJSON serializer.CompClassRequestJSON) (ds.CompClassRequest, error) {
	var compclassrequest ds.CompClassRequest
	err := r.db.Model(&compclassrequest).Where("complex_class_id = ? and big_o_request_id = ?", compclass_id, big_o_request_id).Updates(serializer.CompClassRequestFromJSON(CompClassRequestJSON)).First(&compclassrequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.CompClassRequest{}, fmt.Errorf("%w: класса сложности в заявке", ErrNotFound)
		}
		return ds.CompClassRequest{}, err
	}
	return compclassrequest, nil
}
