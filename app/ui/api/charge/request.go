package charge

type ChargeRequest struct {
	Amount int `json:"amount"`
}

type ChargeCancelRequest struct {
	Amount int `json:"amount"`
}
