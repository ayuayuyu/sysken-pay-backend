package purchase

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"

	"github.com/google/uuid"
)

type CreatePurchaseUseCase interface {
	CreatePurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error)
}

type CreatePurchaseServiceImpl struct {
	purchaseCreateRepo repository.PurchaseRepository
}

func NewCreatePurchaseUseCase(
	purchaseCreateRepo repository.PurchaseRepository,
) *CreatePurchaseServiceImpl {
	return &CreatePurchaseServiceImpl{
		purchaseCreateRepo: purchaseCreateRepo,
	}
}

func (s *CreatePurchaseServiceImpl) CreatePurchase(
	ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error) {

	createdPurchase, err := s.purchaseCreateRepo.CreatePurchase(ctx, userID, items)
	if err != nil {
		return nil, err
	}

	return createdPurchase, nil
}
