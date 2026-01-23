package charge

import (
	"context"
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"

	"github.com/google/uuid"
)

type ChargeCancelUseCase interface {
	ChargeCancel(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error)
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
	ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {

	canceledCharge, err := s.chargeCancelRepo.ChargeCancel(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	return canceledCharge, nil
}
