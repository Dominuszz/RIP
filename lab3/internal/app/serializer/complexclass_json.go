package serializer

import "lab3/internal/app/ds"

type ComplexClassListResponse struct {
	Items      []ComplexClassJSON `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

type ComplexClassJSON struct {
	ID          uint    `json:"compclass_id"`
	IMG         string  `json:"img"`
	Complexity  string  `json:"complexity"`
	Degree      float64 `json:"degree"`
	DegreeText  string  `json:"degree_text"`
	Description string  `json:"description"`
	IsDelete    bool    `json:"is_delete"`
}

func CompClassToJSON(complexclass ds.ComplexClass) ComplexClassJSON {
	return ComplexClassJSON{
		ID:          complexclass.ID,
		IMG:         complexclass.IMG,
		Complexity:  complexclass.Complexity,
		Degree:      complexclass.Degree,
		DegreeText:  complexclass.DegreeText,
		Description: complexclass.Description,
		IsDelete:    complexclass.IsDelete,
	}
}

func CompClassFromJSON(compclassJSON ComplexClassJSON) ds.ComplexClass {
	return ds.ComplexClass{
		IMG:         compclassJSON.IMG,
		Complexity:  compclassJSON.Complexity,
		Degree:      compclassJSON.Degree,
		DegreeText:  compclassJSON.DegreeText,
		Description: compclassJSON.Description,
		IsDelete:    compclassJSON.IsDelete,
	}
}
