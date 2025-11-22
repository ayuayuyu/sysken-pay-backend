package item

import (
	"context"
	"sysken-pay-api/app/domain/object/item"
	"sysken-pay-api/app/domain/repository"
)

//TODO ドメイン層のインターフェースに接続して処理を完成させる

type GetAllItemsUseCase interface {
	GetAllItems(ctx context.Context) ([]*item.Item, error)
}

type GetAllItemsServiceImpl struct {
	itemGetAllRepo repository.ItemRepository
}

func NewGetAllItemsUseCase(
	itemGetAllRepo repository.ItemRepository,
) *GetAllItemsServiceImpl {
	return &GetAllItemsServiceImpl{
		itemGetAllRepo: itemGetAllRepo,
	}
}

func (s *GetAllItemsServiceImpl) GetAllItems(
	ctx context.Context) ([]*item.Item, error) {

	allItems, err := s.itemGetAllRepo.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	return allItems, nil
}
