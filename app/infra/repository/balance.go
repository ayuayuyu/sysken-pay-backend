package repository

import (
	"context"
	"database/sql"
	"time"
	domainbalance "sysken-pay-api/app/domain/object/balance"
	"sysken-pay-api/app/domain/repository"
)

var _ repository.BalanceRepository = (*BalanceRepositoryImpl)(nil)

type BalanceRepositoryImpl struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepositoryImpl {
	return &BalanceRepositoryImpl{db: db}
}

// GetBalance はユーザーの現在の残高を取得する
func (r *BalanceRepositoryImpl) GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error) {
	executor := getExecutor(ctx, r.db)

	var currentBalance int
	query := `SELECT IFNULL(SUM(amount), 0) FROM balance WHERE user_id = ?`
	if err := executor.QueryRowContext(ctx, query, userID).Scan(&currentBalance); err != nil {
		return nil, err
	}

	return domainbalance.NewBalance(userID, currentBalance)
}

// GetPurchaseHistories はユーザーの購入履歴を取得する
func (r *BalanceRepositoryImpl) GetPurchaseHistories(ctx context.Context, userID string, limit, offset int) (
	histories []*domainbalance.PurchaseHistory, totalCount int, totalAmount int, err error,
) {
	executor := getExecutor(ctx, r.db)

	// 全件数と合計支払い金額を取得
	queryAgg := `
		SELECT COUNT(*), IFNULL(SUM(pi.quantity * pi.price_at_purchase), 0)
		FROM purchase p
		INNER JOIN purchase_item pi ON pi.purchase_id = p.id
		WHERE p.user_id = ?
	`
	if err = executor.QueryRowContext(ctx, queryAgg, userID).Scan(&totalCount, &totalAmount); err != nil {
		return nil, 0, 0, err
	}

	// 履歴を取得
	queryRows := `
		SELECT i.name, pi.quantity, pi.price_at_purchase, p.created_at
		FROM purchase p
		INNER JOIN purchase_item pi ON pi.purchase_id = p.id
		INNER JOIN item i ON i.id = pi.item_id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := executor.QueryContext(ctx, queryRows, userID, limit, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			itemName   string
			quantity   int
			price      int
			purchaseAt time.Time
		)
		if err = rows.Scan(&itemName, &quantity, &price, &purchaseAt); err != nil {
			return nil, 0, 0, err
		}
		h, err := domainbalance.NewPurchaseHistory(itemName, quantity, price, purchaseAt)
		if err != nil {
			return nil, 0, 0, err
		}
		histories = append(histories, h)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, 0, err
	}

	return histories, totalCount, totalAmount, nil
}
