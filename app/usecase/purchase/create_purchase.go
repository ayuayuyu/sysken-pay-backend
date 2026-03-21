package purchase

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"
	domainservice "sysken-pay-api/app/domain/service/purchase"
)

type CreatePurchaseUseCase interface {
	CreatePurchase(ctx context.Context, userID string, inputs []PurchaseItemInput) (*purchase.Purchase, error)
}

type CreatePurchaseServiceImpl struct {
	purchaseCreateRepo repository.PurchaseRepository
	itemRepo           repository.ItemRepository
	balanceRepo        repository.BalanceRepository
	txManager          repository.Transaction
}

func NewCreatePurchaseUseCase(
	purchaseCreateRepo repository.PurchaseRepository,
	itemRepo repository.ItemRepository,
	balanceRepo repository.BalanceRepository,
	txManager repository.Transaction,
) *CreatePurchaseServiceImpl {
	return &CreatePurchaseServiceImpl{
		purchaseCreateRepo: purchaseCreateRepo,
		itemRepo:           itemRepo,
		balanceRepo:        balanceRepo,
		txManager:          txManager,
	}
}

func (s *CreatePurchaseServiceImpl) CreatePurchase(
	ctx context.Context, userID string, inputs []PurchaseItemInput) (*purchase.Purchase, error) {

	// ドメインオブジェクト生成 & 合計金額の算出
	items := make([]purchase.PurchaseItem, len(inputs))
	var totalAmount int
	for i, input := range inputs {
		pi, err := purchase.NewPurchaseItem(input.ItemID, input.Quantity)
		if err != nil {
			return nil, err
		}
		items[i] = pi

		it, err := s.itemRepo.GetItemByID(ctx, input.ItemID)
		if err != nil {
			return nil, err
		}
		totalAmount += it.Price() * input.Quantity
	}

	p, err := purchase.NewPurchase(userID, items)
	if err != nil {
		return nil, err
	}

	var createdPurchase *purchase.Purchase
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		// トランザクション内で残高チェック（同時実行安全）
		b, err := s.balanceRepo.GetBalance(ctx, userID)
		if err != nil {
			return err
		}
		if err := domainservice.CanPurchase(b.Balance(), totalAmount); err != nil {
			return err
		}

		createdPurchase, err = s.purchaseCreateRepo.CreatePurchase(ctx, p)
		return err
	})
	if err != nil {
		return nil, err
	}

	return createdPurchase, nil
}
