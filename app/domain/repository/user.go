package repository

import (
	"context"
	"sysken-pay-api/app/domain/object/user"

	"github.com/google/uuid"
)

//TODO　ユーザを登録するインターフェースを作成する
//データベースで必要な入力と出力のインターフェースの作成

type UserRepository interface {
	// ユーザーを新規作成して保存する
	// 保存に成功した場合は保存したユーザーを返す
	InsertUser(ctx context.Context, userName string) (*user.User, error)

	// ユーザー情報を更新して保存する
	// 保存に成功した場合は保存したユーザーを返す
	UpdateUser(ctx context.Context, userID uuid.UUID, userName string) (*user.User, error)
}
