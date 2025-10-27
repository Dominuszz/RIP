package serializer

import "lab3/internal/app/ds"

type CompClassRequestJSON struct {
	ID             uint `json:"comp_class_request_id"`
	BigORequestID  uint `json:"big_o_request_id"`
	ComplexClassID uint `json:"complexclass_id"`
	ArraySize      uint `json:"array_size"`
}

func CompClassRequestToJSON(compclassrequest ds.CompClassRequest) CompClassRequestJSON {
	return CompClassRequestJSON{
		ID:             compclassrequest.ID,
		BigORequestID:  compclassrequest.BigORequestID,
		ComplexClassID: compclassrequest.ComplexClassID,
		ArraySize:      compclassrequest.ArraySize,
	}
}

func CompClassRequestFromJSON(CompClassRequestJSON CompClassRequestJSON) ds.CompClassRequest {
	return ds.CompClassRequest{
		ArraySize: CompClassRequestJSON.ArraySize,
	}
}
