package balance

import (
	"context"
	domainbalance "sysken-pay-api/app/domain/object/balance"
	"sysken-pay-api/app/domain/repository"
)

type GetBalanceUseCase interface {
	GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error)
}

type GetBalanceServiceImpl struct {
	balanceRepo repository.BalanceRepository
}

func NewGetBalanceUseCase(balanceRepo repository.BalanceRepository) *GetBalanceServiceImpl {
	return &GetBalanceServiceImpl{balanceRepo: balanceRepo}
}

func (s *GetBalanceServiceImpl) GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error) {
	return s.balanceRepo.GetBalance(ctx, userID)
}
