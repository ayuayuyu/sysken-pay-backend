package repository

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"

	"github.com/google/uuid"
)

type PurchaseRepository interface {
	// 購入を作成する
	// 購入に成功した場合は購入情報を返す
	CreatePurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error)

	// 購入をキャンセルする
	// キャンセルに成功した場合はキャンセルした購入情報を返す
	CancelPurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error)
}
