package charge

import (
	"context"
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"
)

type ChargeAmountUseCase interface {
	ChargeAmount(ctx context.Context, userID string, amount int) (*charge.Charge, error)
}

type ChargeAmountServiceImpl struct {
	chargeAmountRepo repository.ChargeRepository
}

func NewChargeAmountUseCase(
	chargeAmountRepo repository.ChargeRepository,
) *ChargeAmountServiceImpl {
	return &ChargeAmountServiceImpl{
		chargeAmountRepo: chargeAmountRepo,
	}
}

func (s *ChargeAmountServiceImpl) ChargeAmount(
	ctx context.Context, userID string, amount int) (*charge.Charge, error) {

	c, err := charge.NewCharge(userID, amount)
	if err != nil {
		return nil, err
	}

	chargedAmount, err := s.chargeAmountRepo.ChargeAmount(ctx, c)
	if err != nil {
		return nil, err
	}

	return chargedAmount, nil
}
