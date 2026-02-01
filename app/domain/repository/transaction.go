package repository

import (
	"context"
)

// Transaction はトランザクション制御を行うインターフェースです。
type Transaction interface {
	// Do はトランザクション内で関数 f を実行します。
	Do(ctx context.Context, f func(ctx context.Context) error) error
}
