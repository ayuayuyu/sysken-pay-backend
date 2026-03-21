package purchase

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"
)

type CreatePurchaseUseCase interface {
	CreatePurchase(ctx context.Context, userID string, items []purchase.PurchaseItem) (*purchase.Purchase, error)
}

type CreatePurchaseServiceImpl struct {
	purchaseCreateRepo repository.PurchaseRepository
	txManager          repository.Transaction
}

func NewCreatePurchaseUseCase(
	purchaseCreateRepo repository.PurchaseRepository,
	txManager repository.Transaction,
) *CreatePurchaseServiceImpl {
	return &CreatePurchaseServiceImpl{
		purchaseCreateRepo: purchaseCreateRepo,
		txManager:          txManager,
	}
}

func (s *CreatePurchaseServiceImpl) CreatePurchase(
	ctx context.Context, userID string, items []purchase.PurchaseItem) (*purchase.Purchase, error) {

	p, err := purchase.NewPurchase(userID, items)
	if err != nil {
		return nil, err
	}

	var createdPurchase *purchase.Purchase

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		var err error
		createdPurchase, err = s.purchaseCreateRepo.CreatePurchase(ctx, p)
		return err
	})

	if err != nil {
		return nil, err
	}

	return createdPurchase, nil
}
