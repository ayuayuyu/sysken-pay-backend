package charge

type ChargeReseponse struct {
	Status     string `json:"status"`
	ChargeID   string `json:"charge_id"`
	Amount     int    `json:"charge_amount"`
	UserID     string `json:"user_id"`
	Balance    int    `json:"balance"`
	CreatedAt  string `json:"created_at"`
}

func toChargeResponse(chargeID string, amount int, userID string, balance int, createdAt string) *ChargeReseponse {
	return &ChargeReseponse{
		Status:    "success",
		ChargeID:  chargeID,
		Amount:    amount,
		UserID:    userID,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}