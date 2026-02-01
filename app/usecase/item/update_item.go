package item

import (
	"context"
	"sysken-pay-api/app/domain/object/item"
	"sysken-pay-api/app/domain/repository"
)

//TODO ドメイン層のインターフェースに接続して処理を完成させる

type UpdateItemUseCase interface {
	UpdateItem(ctx context.Context, janCode string, itemName string, price int) (*item.Item, error)
}

type UpdateItemServiceImpl struct {
	itemUpdateRepo repository.ItemRepository
}

func NewUpdateItemUseCase(
	itemUpdateRepo repository.ItemRepository,
) *UpdateItemServiceImpl {
	return &UpdateItemServiceImpl{
		itemUpdateRepo: itemUpdateRepo,
	}
}

func (s *UpdateItemServiceImpl) UpdateItem(
	ctx context.Context, janCode string, itemName string, price int) (*item.Item, error) {

	i, err := item.UpdateItem(janCode, itemName, price)
	if err != nil {
		return nil, err
	}

	updateItem, err := s.itemUpdateRepo.UpdateItem(ctx, i)
	if err != nil {
		return nil, err
	}

	return updateItem, nil
}
