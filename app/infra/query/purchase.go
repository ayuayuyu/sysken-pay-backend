package query

import (
	"context"
	"database/sql"
	"sysken-pay-api/app/domain/object/purchase"
	"sysken-pay-api/app/domain/repository"
	"time"

	"github.com/google/uuid"
)

var _ repository.PurchaseRepository = (*PurchaseRepositoryImpl)(nil)

type PurchaseRepositoryImpl struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepositoryImpl {
	return &PurchaseRepositoryImpl{
		db: db,
	}
}

func (r *PurchaseRepositoryImpl) CreatePurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error) {
	p, err := purchase.NewPurchase(userID, items)
	if err != nil {
		return nil, err
	}

	// トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 合計金額を計算し、各商品の価格を取得する
	var totalAmount int
	type itemData struct {
		itemID   int
		quantity int
		price    int
	}
	var itemsToInsert []itemData

	for _, pi := range items {
		var price int
		// 商品の価格を取得
		err = tx.QueryRowContext(ctx, "SELECT price FROM item WHERE id = ?", pi.ItemID()).Scan(&price)
		if err != nil {
			return nil, err
		}
		totalAmount += price * pi.Quantity()
		itemsToInsert = append(itemsToInsert, itemData{
			itemID:   pi.ItemID(),
			quantity: pi.Quantity(),
			price:    price,
		})
	}

	// purchaseテーブルへの挿入
	queryPurchase := `INSERT INTO purchase (user_id) VALUES (?)`
	res, err := tx.ExecContext(ctx, queryPurchase, userID)
	if err != nil {
		return nil, err
	}

	purchaseID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// purchase_itemテーブルへの挿入
	queryPurchaseItem := `INSERT INTO purchase_item (purchase_id, item_id, quantity, price_at_purchase) VALUES (?, ?, ?, ?)`
	stmt, err := tx.PrepareContext(ctx, queryPurchaseItem)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, it := range itemsToInsert {
		_, err := stmt.ExecContext(ctx, purchaseID, it.itemID, it.quantity, it.price)
		if err != nil {
			return nil, err
		}
	}

	// balanceテーブルへの挿入 (購入なのでamountはマイナス)
	queryBalance := `INSERT INTO balance (user_id, purchase_id, amount) VALUES (?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryBalance, userID, purchaseID, -totalAmount)
	if err != nil {
		return nil, err
	}

	p.SetID(int(purchaseID))
	p.SetCreatedAt(time.Now())

	return p, nil
}

func (r *PurchaseRepositoryImpl) CancelPurchase(ctx context.Context, userID uuid.UUID, items []purchase.PurchaseItem) (*purchase.Purchase, error) {
	p, err := purchase.DeletePurchase(userID, items)
	if err != nil {
		return nil, err
	}

	// トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 合計金額を計算
	var totalAmount int

	for _, pi := range items {
		var price int
		// 商品の価格を取得
		err = tx.QueryRowContext(ctx, "SELECT price FROM item WHERE id = ?", pi.ItemID()).Scan(&price)
		if err != nil {
			return nil, err
		}
		totalAmount += price * pi.Quantity()
	}

	// balanceテーブルへの挿入 (キャンセルなのでamountはプラス)
	queryBalance := `INSERT INTO balance (user_id, amount) VALUES (?, ?)`
	_, err = tx.ExecContext(ctx, queryBalance, userID, totalAmount)
	if err != nil {
		return nil, err
	}

	return p, nil
}
