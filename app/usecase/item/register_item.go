package item

import (
	"context"
	"sysken-pay-api/app/domain/object/item"
	"sysken-pay-api/app/domain/repository"
)

//TODO ドメイン層のインターフェースに接続して処理を完成させる

type RegisterItemUseCase interface {
	RegisterItem(ctx context.Context, janCode string, itemName string, price int) (*item.Item, error)
}

type RegisterItemServiceImpl struct {
	itemRegisterRepo repository.ItemRepository
}

func NewRegisterItemUseCase(
	itemRegisterRepo repository.ItemRepository,
) *RegisterItemServiceImpl {
	return &RegisterItemServiceImpl{
		itemRegisterRepo: itemRegisterRepo,
	}
}

func (s *RegisterItemServiceImpl) RegisterItem(
	ctx context.Context, janCode string, itemName string, price int) (*item.Item, error) {

	createdItem, err := s.itemRegisterRepo.InsertItem(ctx, janCode, itemName, price)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}
