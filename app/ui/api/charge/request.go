package charge

type PostChargeRequest struct {
	Amount int `json:"amount"`
}

type PostChargeCancelRequest struct {
	Amount int `json:"amount"`
}
