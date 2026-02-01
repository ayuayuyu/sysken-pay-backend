package repository

import (
	"context"
	"database/sql"
	"sysken-pay-api/app/infra/transaction"
)

// DBExecutor は *sql.DB と *sql.Tx の共通メソッドを持つインターフェースです
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func getExecutor(ctx context.Context, db *sql.DB) DBExecutor {
	if tx, ok := transaction.GetTx(ctx); ok {
		return tx
	}
	return db
}
