package serializer

type CalculationRequestJSON struct {
	RequestID   int                 `json:"request_id"`
	CompClasses []CompClassCalcJSON `json:"compclasses"`
}

type CompClassCalcJSON struct {
	Complexity string  `json:"complexity"`
	Degree     float64 `json:"degree"`
	ArraySize  uint    `json:"array_size"`
}

type CalculationResultJSON struct {
	RequestID            int     `json:"request_id"`
	CalculatedTime       float64 `json:"calculated_time"`
	CalculatedComplexity string  `json:"calculated_complexity"`
	Success              bool    `json:"success"`
	Status               string  `json:"status"`
	Message              string  `json:"message"`
}
