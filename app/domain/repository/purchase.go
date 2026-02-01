package repository

import (
	"context"
	"sysken-pay-api/app/domain/object/purchase"
)

type PurchaseRepository interface {
	// 購入を作成する
	// 購入に成功した場合は購入情報を返す
	CreatePurchase(ctx context.Context, p *purchase.Purchase) (*purchase.Purchase, error)

	// 購入をキャンセルする
	// キャンセルに成功した場合はキャンセルした購入情報を返す
	CancelPurchase(ctx context.Context, p *purchase.Purchase) (*purchase.Purchase, error)
}
