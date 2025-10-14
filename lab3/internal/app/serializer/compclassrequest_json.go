package serializer

import "lab3/internal/app/ds"

type CompClassRequestJSON struct {
	ID             uint    `json:"comp_class_request_id"`
	BigORequestID  uint    `json:"big_o_request_id"`
	ComplexClassID uint    `json:"complexclass_id"`
	Complexity     string  `json:"complexity"`
	Degree         float64 `json:"degree"`
	DegreeText     string  `json:"degree_text"`
	ArraySize      uint    `json:"array_size"`
}

func CompClassRequestToJSON(compclassrequest ds.CompClassRequest) CompClassRequestJSON {
	return CompClassRequestJSON{
		ID:             compclassrequest.ID,
		BigORequestID:  compclassrequest.BigORequestID,
		ComplexClassID: compclassrequest.ComplexClassID,
		Complexity:     compclassrequest.Complexity,
		Degree:         compclassrequest.Degree,
		DegreeText:     compclassrequest.DegreeText,
		ArraySize:      compclassrequest.ArraySize,
	}
}

func CompClassRequestFromJSON(CompClassRequestJSON CompClassRequestJSON) ds.CompClassRequest {
	return ds.CompClassRequest{
		Complexity: CompClassRequestJSON.Complexity,
		Degree:     CompClassRequestJSON.Degree,
		DegreeText: CompClassRequestJSON.DegreeText,
		ArraySize:  CompClassRequestJSON.ArraySize,
	}
}
