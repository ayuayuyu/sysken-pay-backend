package charge

import (
	"sysken-pay-api/app/domain/object/charge"
)

type ChargeResponse struct {
	Status    string `json:"status"`
	ChargeID  int    `json:"charge_id"`
	Amount    int    `json:"charge_amount"`
	UserID    string `json:"user_id"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

func toPostChargeResponse(charge *charge.Charge) *ChargeResponse {
	return &ChargeResponse{
		Status:    "success",
		ChargeID:  charge.ID(),
		Amount:    charge.Amount(),
		UserID:    charge.UserID(),
		Balance:   charge.Balance(),
		CreatedAt: charge.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
	}
}

type ChargeCancelResponse struct {
	Status    string `json:"status"`
	ChargeID  int    `json:"charge_id"`
	Amount    int    `json:"canceled_amount"`
	UserID    string `json:"user_id"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
}

func toPostChargeCancelResponse(charge *charge.Charge) *ChargeCancelResponse {
	return &ChargeCancelResponse{
		Status:    "success",
		ChargeID:  charge.ID(),
		Amount:    charge.Amount(),
		UserID:    charge.UserID(),
		Balance:   charge.Balance(),
		CreatedAt: charge.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
	}
}
