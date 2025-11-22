package item

import (
	"context"
	"sysken-pay-api/app/domain/object/item"
	"sysken-pay-api/app/domain/repository"
)

//TODO ドメイン層のインターフェースに接続して処理を完成させる

type FindItemByJanCodeUseCase interface {
	GetItemByJanCode(ctx context.Context, janCode string) (*item.Item, error)
}

type FindItemByJanCodeServiceImpl struct {
	itemFindRepo repository.ItemRepository
}

func NewFindItemByJanCodeUseCase(
	itemFindByJanCodeRepo repository.ItemRepository,
) *FindItemByJanCodeServiceImpl {
	return &FindItemByJanCodeServiceImpl{
		itemFindRepo: itemFindByJanCodeRepo,
	}
}

func (s *FindItemByJanCodeServiceImpl) GetItemByJanCode(
	ctx context.Context, janCode string) (*item.Item, error) {

	foundItemByJanCode, err := s.itemFindRepo.GetItemByJanCode(ctx, janCode)
	if err != nil {
		return nil, err
	}

	return foundItemByJanCode, nil
}
