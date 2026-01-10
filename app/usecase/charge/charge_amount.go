package charge

import (
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"
)

type ChargeAmountUseCase interface {
	ChargeAmount(userID string, amount int) (*charge.Charge, error)
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
	userID string, amount int) (*charge.Charge, error) {

	chargedAmount, err := s.chargeAmountRepo.ChargeAmount(userID, amount)
	if err != nil {
		return nil, err
	}

	return chargedAmount, nil
}
