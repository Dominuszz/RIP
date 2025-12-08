package serializer

import (
	"encoding/json"
	"lab3/internal/app/ds"
	"time"
)

type RuTime time.Time

func (t RuTime) MarshalJSON() ([]byte, error) {
	// Формат: 31.12.2025 23:59:59
	formatted := time.Time(t).Format("02.01.2006 15:04:05")
	return json.Marshal(formatted)
}

func (t *RuTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	parsed, err := time.Parse("02.01.2006 15:04:05", str)
	if err != nil {
		return err
	}
	*t = RuTime(parsed)
	return nil
}

// Модифицированная структура для JSON
type BigORequestJSON struct {
	ID                   uint    `json:"bigo_request_id"`
	Status               string  `json:"status"`
	DateCreate           RuTime  `json:"date_create"`
	Creator_Login        string  `json:"creator_login"`
	Moderator_Login      *string `json:"moderator_login"`
	DateUpdate           *RuTime `json:"date_update"`
	DateFinish           *RuTime `json:"date_finish"`
	CalculatedTime       float64 `json:"calculated_time"`
	CalculatedComplexity string  `json:"calculated_complexity"`
}

func BigORequestToJSON(bigorequest ds.BigORequest, creator_login string, moderator_login string) BigORequestJSON {
	var upd_date, fin_date *RuTime
	if bigorequest.DateUpdate.Valid {
		tmp := RuTime(bigorequest.DateUpdate.Time)
		upd_date = &tmp
	}
	if bigorequest.DateFinish.Valid {
		tmp := RuTime(bigorequest.DateFinish.Time)
		fin_date = &tmp
	}
	var m_login *string
	if moderator_login != "" {
		m_login = &moderator_login
	}

	return BigORequestJSON{
		ID:                   bigorequest.ID,
		Status:               bigorequest.Status,
		DateCreate:           RuTime(bigorequest.DateCreate),
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
type ResultUpdateJSON struct {
	CalculatedComplexity string  `json:"calculated_complexity"`
	CalculatedTime       float64 `json:"calculated_time"`
	AuthKey              string  `json:"auth_key"`
}
type AsyncCalculationJSON struct {
	PK int `json:"pk"`
}
