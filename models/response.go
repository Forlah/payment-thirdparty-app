package models

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type PaymentResponsePayload struct {
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}
