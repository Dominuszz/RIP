package serializer

import "lab3/internal/app/ds"

type ComplexClassJSON struct {
	ID          uint    `json:"compclass_id"`
	Complexity  string  `json:"complexity"`
	Degree      float64 `json:"degree"`
	DegreeText  string  `json:"degree_text"`
	Description string  `json:"description"`
	IsDelete    bool    `json:"is_delete"`
}

func CompClassToJSON(complexclass ds.ComplexClass) ComplexClassJSON {
	return ComplexClassJSON{
		ID:          complexclass.ID,
		Complexity:  complexclass.Complexity,
		Degree:      complexclass.Degree,
		DegreeText:  complexclass.DegreeText,
		Description: complexclass.Description,
		IsDelete:    complexclass.IsDelete,
	}
}

func CompClassFromJSON(compclassJSON ComplexClassJSON) ds.ComplexClass {
	return ds.ComplexClass{
		Complexity:  compclassJSON.Complexity,
		Degree:      compclassJSON.Degree,
		DegreeText:  compclassJSON.DegreeText,
		Description: compclassJSON.Description,
		IsDelete:    compclassJSON.IsDelete,
	}
}
