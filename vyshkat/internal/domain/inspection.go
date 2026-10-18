package domain

// Лежит в domain чтобы fleet и inspection не зависели друг от друга
type InspectionResult struct {
	Approved bool
	Reason   string
}

func Approved() InspectionResult {
	return InspectionResult{Approved: true, Reason: "техосмотр пройден"}
}

func Rejected(reason string) InspectionResult {
	return InspectionResult{Approved: false, Reason: reason}
}
