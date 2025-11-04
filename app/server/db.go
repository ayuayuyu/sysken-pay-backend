package server

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// NewDB は、mysql.Config を使って *sql.DB を初期化します。
func NewDB(config *mysql.Config) (*sql.DB, error) {
	driverName := "mysql"
	dsn := config.FormatDSN()

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	// 実際に接続確認を行う（設定ミスを早期検出）
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping failed: %w", err)
	}

	return db, nil
}
