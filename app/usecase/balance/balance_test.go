package balance

import (
	"context"
	"errors"
	"testing"
	"time"

	domainbalance "sysken-pay-api/app/domain/object/balance"
)

// --- mock ---

type mockBalanceRepo struct {
	getBalanceFunc         func(ctx context.Context, userID string) (*domainbalance.Balance, error)
	getPurchaseHistoryFunc func(ctx context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error)
}

func (m *mockBalanceRepo) GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error) {
	return m.getBalanceFunc(ctx, userID)
}

func (m *mockBalanceRepo) GetPurchaseHistories(ctx context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
	return m.getPurchaseHistoryFunc(ctx, userID, limit, offset)
}

// helper
func newTestHistory(itemName string, quantity, price int) *domainbalance.PurchaseHistory {
	h, _ := domainbalance.NewPurchaseHistory(itemName, quantity, price, time.Now())
	return h
}

// --- GetBalance ---

func TestGetBalance_Success(t *testing.T) {
	repo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return domainbalance.NewBalance(userID, 1500)
		},
	}
	uc := NewGetBalanceUseCase(repo)
	result, err := uc.GetBalance(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetBalance should succeed: %v", err)
	}
	if result.Balance() != 1500 {
		t.Errorf("Balance() = %d, want 1500", result.Balance())
	}
}

func TestGetBalance_RepoError(t *testing.T) {
	repo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewGetBalanceUseCase(repo)
	if _, err := uc.GetBalance(context.Background(), "user-1"); err == nil {
		t.Error("GetBalance should propagate repo error")
	}
}

// --- GetPurchaseHistories ---

func TestGetPurchaseHistories_Success(t *testing.T) {
	histories := []*domainbalance.PurchaseHistory{
		newTestHistory("コーラ", 2, 150),
		newTestHistory("パン", 1, 200),
	}
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			return histories, 2, 500, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	result, totalAmount, pageInfo, err := uc.GetPurchaseHistories(context.Background(), "user-1", 1, 20)
	if err != nil {
		t.Fatalf("GetPurchaseHistories should succeed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}
	if totalAmount != 500 {
		t.Errorf("totalAmount = %d, want 500", totalAmount)
	}
	if pageInfo == nil {
		t.Fatal("pageInfo should not be nil")
	}
	if pageInfo.TotalPage() != 1 {
		t.Errorf("TotalPage() = %d, want 1", pageInfo.TotalPage())
	}
}

func TestGetPurchaseHistories_Pagination_PrevAndNext(t *testing.T) {
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			// limit=5, offset=5 (page=2, perPage=5), totalCount=15 → 3 pages
			return []*domainbalance.PurchaseHistory{newTestHistory("item", 1, 100)}, 15, 100, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	_, _, pageInfo, err := uc.GetPurchaseHistories(context.Background(), "user-1", 2, 5)
	if err != nil {
		t.Fatalf("GetPurchaseHistories should succeed: %v", err)
	}
	if pageInfo.TotalPage() != 3 {
		t.Errorf("TotalPage() = %d, want 3", pageInfo.TotalPage())
	}
	if pageInfo.PrevPage() == nil || *pageInfo.PrevPage() != 1 {
		t.Error("PrevPage should be 1")
	}
	if pageInfo.NextPage() == nil || *pageInfo.NextPage() != 3 {
		t.Error("NextPage should be 3")
	}
}

func TestGetPurchaseHistories_FirstPage_NoPrev(t *testing.T) {
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			return []*domainbalance.PurchaseHistory{newTestHistory("item", 1, 100)}, 25, 100, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	_, _, pageInfo, err := uc.GetPurchaseHistories(context.Background(), "user-1", 1, 10)
	if err != nil {
		t.Fatalf("GetPurchaseHistories should succeed: %v", err)
	}
	if pageInfo.PrevPage() != nil {
		t.Error("PrevPage should be nil on first page")
	}
	if pageInfo.NextPage() == nil || *pageInfo.NextPage() != 2 {
		t.Error("NextPage should be 2")
	}
}

func TestGetPurchaseHistories_LastPage_NoNext(t *testing.T) {
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			return []*domainbalance.PurchaseHistory{newTestHistory("item", 1, 100)}, 10, 100, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	// page=2, perPage=5, totalCount=10 → totalPage=2, last page
	_, _, pageInfo, err := uc.GetPurchaseHistories(context.Background(), "user-1", 2, 5)
	if err != nil {
		t.Fatalf("GetPurchaseHistories should succeed: %v", err)
	}
	if pageInfo.NextPage() != nil {
		t.Error("NextPage should be nil on last page")
	}
	if pageInfo.PrevPage() == nil || *pageInfo.PrevPage() != 1 {
		t.Error("PrevPage should be 1")
	}
}

func TestGetPurchaseHistories_ZeroCount_TotalPageIsOne(t *testing.T) {
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			return []*domainbalance.PurchaseHistory{}, 0, 0, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	_, _, pageInfo, err := uc.GetPurchaseHistories(context.Background(), "user-1", 1, 20)
	if err != nil {
		t.Fatalf("GetPurchaseHistories should succeed with 0 results: %v", err)
	}
	if pageInfo.TotalPage() != 1 {
		t.Errorf("TotalPage() = %d, want 1 when totalCount=0", pageInfo.TotalPage())
	}
}

func TestGetPurchaseHistories_InvalidPage(t *testing.T) {
	repo := &mockBalanceRepo{}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	if _, _, _, err := uc.GetPurchaseHistories(context.Background(), "user-1", 0, 20); err == nil {
		t.Error("GetPurchaseHistories page=0 should fail")
	}
}

func TestGetPurchaseHistories_InvalidPerPage(t *testing.T) {
	repo := &mockBalanceRepo{}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	if _, _, _, err := uc.GetPurchaseHistories(context.Background(), "user-1", 1, 0); err == nil {
		t.Error("GetPurchaseHistories perPage=0 should fail")
	}
}

func TestGetPurchaseHistories_RepoError(t *testing.T) {
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			return nil, 0, 0, errors.New("db error")
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	if _, _, _, err := uc.GetPurchaseHistories(context.Background(), "user-1", 1, 20); err == nil {
		t.Error("GetPurchaseHistories should propagate repo error")
	}
}

func TestGetPurchaseHistories_CorrectLimitOffset(t *testing.T) {
	var capturedLimit, capturedOffset int
	repo := &mockBalanceRepo{
		getPurchaseHistoryFunc: func(_ context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
			capturedLimit = limit
			capturedOffset = offset
			return []*domainbalance.PurchaseHistory{}, 0, 0, nil
		},
	}
	uc := NewGetPurchaseHistoriesUseCase(repo)
	uc.GetPurchaseHistories(context.Background(), "user-1", 3, 10) // page=3, perPage=10 → offset=20
	if capturedLimit != 10 {
		t.Errorf("limit = %d, want 10", capturedLimit)
	}
	if capturedOffset != 20 {
		t.Errorf("offset = %d, want 20", capturedOffset)
	}
}
