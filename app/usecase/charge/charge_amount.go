package charge

import (
	"context"
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"

	"github.com/google/uuid"
)

type ChargeAmountUseCase interface {
	ChargeAmount(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error)
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
	ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {

	chargedAmount, err := s.chargeAmountRepo.ChargeAmount(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	return chargedAmount, nil
}
