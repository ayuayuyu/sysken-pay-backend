package repository

import (
	"context"
	domainbalance "sysken-pay-api/app/domain/object/balance"
)

type BalanceRepository interface {
	// GetBalance はユーザーの現在の残高を取得する
	GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error)

	// GetPurchaseHistories はユーザーの購入履歴をページネーションで取得する
	// totalCount は全件数、totalAmount は合計支払い金額を返す
	GetPurchaseHistories(ctx context.Context, userID string, page, perPage int) (histories []*domainbalance.PurchaseHistory, totalCount int, totalAmount int, err error)
}
