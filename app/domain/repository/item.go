package repository

import (
	"context"
	"sysken-pay-api/app/domain/object/item"
)

//TODO　商品を登録するインターフェースを作成する
//データベースで必要な入力と出力のインターフェースの作成
//TODO 商品を取得するインターフェースを作成する

type ItemRepository interface {
	// 商品を新規作成して保存する
	// 保存に成功した場合は保存した商品を返す
	InsertItem(ctx context.Context, janCode string, name string, price int) (*item.Item, error)

	// 商品を更新する
	// 更新に成功した場合は更新した商品を返す
	UpdateItem(ctx context.Context, janCode string, name string, price int) (*item.Item, error)

	// 商品をJANコードで取得する
	// 取得に成功した場合は取得した商品を返す
	GetItemByJanCode(ctx context.Context, janCode string) (*item.Item, error)

	//商品を全て取得する
	GetAllItems(ctx context.Context) ([]*item.Item, error)
}
