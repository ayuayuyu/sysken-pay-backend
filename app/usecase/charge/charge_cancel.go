package charge

import (
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"
)

type ChargeCancelUseCase interface {
	ChargeCancel(userID string, amount int) (*charge.Charge, error)
}

type ChargeCancelServiceImpl struct {
	chargeCancelRepo repository.ChargeRepository
}

func NewChargeCancelUseCase(
	chargeCancelRepo repository.ChargeRepository,
) *ChargeCancelServiceImpl {
	return &ChargeCancelServiceImpl{
		chargeCancelRepo: chargeCancelRepo,
	}
}

func (s *ChargeCancelServiceImpl) ChargeCancel(
	userID string, amount int) (*charge.Charge, error) {

	canceledCharge, err := s.chargeCancelRepo.ChargeCancel(userID, amount)
	if err != nil {
		return nil, err
	}

	return canceledCharge, nil
}
