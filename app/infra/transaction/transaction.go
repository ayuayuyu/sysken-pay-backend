package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"sysken-pay-api/app/domain/repository"
)

// key はcontextのキーとして使用する型です
type key struct{}

type transaction struct {
	db *sql.DB
}

// NewTransaction は Transaction の新しいインスタンスを生成します。
func NewTransaction(db *sql.DB) repository.Transaction {
	return &transaction{db: db}
}

// GetTx はコンテキストからトランザクションを取得します。
func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(key{}).(*sql.Tx)
	return tx, ok
}

// Do はトランザクション内で関数 f を実行します。
// f がエラーを返した場合はロールバックし、nil を返した場合はコミットします。
func (t *transaction) Do(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction begin error: %w", err)
	}

	// コンテキストにトランザクションを埋め込む
	ctx = context.WithValue(ctx, key{}, tx)

	// パニック時のリカバリ
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := f(ctx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback error: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error: %w", err)
	}

	return nil
}

// DBConnection は *sql.DB と *sql.Tx の共通メソッドを持つインターフェースです。
type DBConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// GetDB はコンテキストからトランザクションを取得し、存在すればそれを返します。
// トランザクションが存在しない場合は、引数で渡された db (*sql.DB) を返します。
// これにより、リポジトリ層はトランザクションの有無を意識せずに実装できます。
func GetDB(ctx context.Context, db *sql.DB) DBConnection {
	v := ctx.Value(key{})
	if tx, ok := v.(*sql.Tx); ok {
		return tx
	}
	return db
}
