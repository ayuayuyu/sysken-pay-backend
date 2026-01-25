package query

import (
	"context"
	"database/sql"
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"
	"time"

	"github.com/google/uuid"
)

var _ repository.ChargeRepository = (*ChargeRepositoryImpl)(nil)

type ChargeRepositoryImpl struct {
	db *sql.DB
}

func NewChargeRepository(db *sql.DB) *ChargeRepositoryImpl {
	return &ChargeRepositoryImpl{
		db: db,
	}
}

// ChargeAmount は指定された金額をチャージし、チャージ情報を返す
func (r *ChargeRepositoryImpl) ChargeAmount(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {
	c, err := charge.NewCharge(userID, amount)
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

	// 1. charge テーブルへ挿入
	queryCharge := `INSERT INTO charge (user_id, amount, created_at) VALUES (?, ?, ?)`
	createdAt := time.Now()
	res, err := tx.ExecContext(ctx, queryCharge, c.UserID(), c.Amount(), createdAt)
	if err != nil {
		return nil, err
	}

	chargeID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	c.SetID(int(chargeID))
	c.SetCreatedAt(createdAt)

	// 2. balance テーブルへ挿入 (残高追加)
	queryBalance := `INSERT INTO balance (user_id, charge_id, amount, created_at) VALUES (?, ?, ?, ?)`
	// チャージなので amount は正の値
	if _, err := tx.ExecContext(ctx, queryBalance, c.UserID(), int(chargeID), c.Amount(), createdAt); err != nil {
		return nil, err
	}

	// 3. 現在の残高を集計して取得
	var currentBalance int
	querySum := `SELECT IFNULL(SUM(amount), 0) FROM balance WHERE user_id = ?`
	if err := tx.QueryRowContext(ctx, querySum, c.UserID()).Scan(&currentBalance); err != nil {
		return nil, err
	}

	c.SetBalance(currentBalance)

	return c, nil
}


// ChargeCancel は本来「特定のチャージIDを取り消す」べきですが、インターフェースが (userID, amount) なので
// 「指定された金額分のチャージを取り消す（=残高を減らす）」処理として実装します。
// chargeテーブルのレコードを論理削除し、balanceテーブルでマイナスを記録します。
func (r *ChargeRepositoryImpl) ChargeCancel(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {
	// 一旦ダミーのチャージオブジェクトを作成（返り値用）
	c := &charge.Charge{}
	c.SetUserID(userID)
	// 本来は負の値を入れたいが、domain objectの制約で入れられない場合がある
	// ここでは amount そのものは正 (キャンセル額) として扱う

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

	// 1. 取り消し対象のチャージを検索 (まだ削除されていない、金額が一致する最新のものなど)
	// ※厳密な仕様が不明なため、ここでは「即座に論理削除されたチャージレコード」を作成する方式をとるか、
	// あるいは既存のレコードを探して更新するか。
	// 既存コードの流れに従い、新規作成→即時論理削除を行うパターンで実装します。
	// (ただし、chargeテーブルはamount>0制約があるため、マイナスレコードは作れません。
	//  キャンセル履歴として プラスの額で作成し、即座に deleted_at を入れて無効扱いにするのが妥当です)

	queryCharge := `INSERT INTO charge (user_id, amount, created_at, deleted_at) VALUES (?, ?, ?, ?)`
	now := time.Now()
	res, err := tx.ExecContext(ctx, queryCharge, userID, amount, now, now)
	if err != nil {
		return nil, err
	}

	chargeID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// 返り値設定
	c.SetID(int(chargeID))
	c.SetCreatedAt(now)
	// cancelなのでdeletedAtもセット
	c.SetDeletedAt(now)

	// 2. balance テーブルへ挿入 (残高減少)
	queryBalance := `INSERT INTO balance (user_id, charge_id, amount, created_at) VALUES (?, ?, ?, ?)`
	// キャンセルなので amount は負の値にする
	if _, err := tx.ExecContext(ctx, queryBalance, userID, int(chargeID), -amount, now); err != nil {
		return nil, err
	}

	// 3. 現在の残高を集計して取得
	var currentBalance int
	querySum := `SELECT IFNULL(SUM(amount), 0) FROM balance WHERE user_id = ?`
	if err := tx.QueryRowContext(ctx, querySum, userID).Scan(&currentBalance); err != nil {
		return nil, err
	}

	c.SetBalance(currentBalance)

	return c, nil
}
