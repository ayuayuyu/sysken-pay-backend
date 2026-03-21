package repository

import (
	"context"
	domainbalance "sysken-pay-api/app/domain/object/balance"
)

type BalanceRepository interface {
	// GetBalance はユーザーの現在の残高を取得する
	GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error)

	// GetPurchaseHistories はユーザーの購入履歴を取得する
	// limit: 取得件数、offset: スキップ件数
	// totalCount は全件数、totalAmount は合計支払い金額を返す
	GetPurchaseHistories(ctx context.Context, userID string, limit, offset int) (histories []*domainbalance.PurchaseHistory, totalCount int, totalAmount int, err error)
}
