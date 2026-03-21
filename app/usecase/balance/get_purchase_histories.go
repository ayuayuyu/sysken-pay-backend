package balance

import (
	"context"
	"errors"
	domainbalance "sysken-pay-api/app/domain/object/balance"
	"sysken-pay-api/app/domain/repository"
)

type GetPurchaseHistoriesUseCase interface {
	GetPurchaseHistories(ctx context.Context, userID string, page, perPage int) (
		[]*domainbalance.PurchaseHistory,
		int,
		*domainbalance.HistoryPage,
		error,
	)
}

type GetPurchaseHistoriesServiceImpl struct {
	balanceRepo repository.BalanceRepository
}

func NewGetPurchaseHistoriesUseCase(balanceRepo repository.BalanceRepository) *GetPurchaseHistoriesServiceImpl {
	return &GetPurchaseHistoriesServiceImpl{balanceRepo: balanceRepo}
}

func (s *GetPurchaseHistoriesServiceImpl) GetPurchaseHistories(
	ctx context.Context, userID string, page, perPage int,
) (histories []*domainbalance.PurchaseHistory, totalAmount int, pageInfo *domainbalance.HistoryPage, err error) {
	if page < 1 {
		return nil, 0, nil, errors.New("page must be >= 1")
	}
	if perPage < 1 {
		return nil, 0, nil, errors.New("per_page must be >= 1")
	}

	histories, totalCount, totalAmount, err := s.balanceRepo.GetPurchaseHistories(ctx, userID, page, perPage)
	if err != nil {
		return nil, 0, nil, err
	}

	totalPage := 1
	if totalCount > 0 {
		totalPage = (totalCount + perPage - 1) / perPage
	}

	var prevPage *int
	if page > 1 {
		p := page - 1
		prevPage = &p
	}
	var nextPage *int
	if page < totalPage {
		p := page + 1
		nextPage = &p
	}

	pageInfo, err = domainbalance.NewHistoryPage(prevPage, nextPage, totalPage)
	if err != nil {
		return nil, 0, nil, err
	}

	return histories, totalAmount, pageInfo, nil
}
