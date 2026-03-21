package purchase

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"
)

type CancelPurchaseUseCase interface {
	CancelPurchase(ctx context.Context, userID string, inputs []PurchaseItemInput) (*purchase.Purchase, error)
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
	ctx context.Context, userID string, inputs []PurchaseItemInput) (*purchase.Purchase, error) {

	items := make([]purchase.PurchaseItem, len(inputs))
	for i, input := range inputs {
		pi, err := purchase.NewPurchaseItem(input.ItemID, input.Quantity)
		if err != nil {
			return nil, err
		}
		items[i] = pi
	}

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
