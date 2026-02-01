package purchase

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"

	"github.com/google/uuid"
)

type CancelPurchaseUseCase interface {
	CancelPurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error)
}

type CancelPurchaseServiceImpl struct {
	purchaseCancelRepo repository.PurchaseRepository
}

func NewCancelPurchaseUseCase(
	purchaseCancelRepo repository.PurchaseRepository,
) *CancelPurchaseServiceImpl {
	return &CancelPurchaseServiceImpl{
		purchaseCancelRepo: purchaseCancelRepo,
	}
}

func (s *CancelPurchaseServiceImpl) CancelPurchase(
	ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error) {

	p, err := purchase.DeletePurchase(userID, items)
	if err != nil {
		return nil, err
	}

	canceledPurchase, err := s.purchaseCancelRepo.CancelPurchase(ctx, p)
	if err != nil {
		return nil, err
	}

	return canceledPurchase, nil
}
