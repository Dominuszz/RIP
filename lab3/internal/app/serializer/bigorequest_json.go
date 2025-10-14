package serializer

import (
	"lab3/internal/app/ds"
	"time"
)

type BigORequestJSON struct {
	ID                   uint       `json:"bigo_request_id"`
	Status               string     `json:"status"`
	DateCreate           time.Time  `json:"date_create"`
	Creator_Login        string     `json:"creator_login"`
	Moderator_Login      *string    `json:"moderator_login"`
	DateUpdate           *time.Time `json:"date_update"`
	DateFinish           *time.Time `json:"date_finish"`
	CalculatedTime       float64    `json:"calculated_time"`
	CalculatedComplexity string     `json:"calculated_complexity"`
}

func BigORequestToJSON(bigorequest ds.BigORequest, creator_login string, moderator_login string) BigORequestJSON {
	var upd_date, fin_date *time.Time
	if bigorequest.DateUpdate.Valid {
		upd_date = &bigorequest.DateUpdate.Time
	}
	if bigorequest.DateFinish.Valid {
		fin_date = &bigorequest.DateFinish.Time
	}
	var m_login *string
	if moderator_login != "" {
		m_login = &moderator_login
	}

	return BigORequestJSON{
		ID:                   bigorequest.ID,
		Status:               bigorequest.Status,
		DateCreate:           bigorequest.DateCreate,
		Creator_Login:        creator_login,
		Moderator_Login:      m_login,
		DateUpdate:           upd_date,
		DateFinish:           fin_date,
		CalculatedTime:       bigorequest.CalculatedTime,
		CalculatedComplexity: bigorequest.CalculatedComplexity,
	}
}

func BigORequestFromJSON(bigorequest BigORequestJSON) ds.BigORequest {
	if bigorequest.CalculatedTime == 0 {
		return ds.BigORequest{}
	}
	return ds.BigORequest{
		CalculatedTime: bigorequest.CalculatedTime,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
