package models

type PostPaymentRequestPayload struct {
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}
