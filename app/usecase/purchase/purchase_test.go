package purchase

import (
	"context"
	"errors"
	"testing"

	domainbalance "sysken-pay-api/app/domain/object/balance"
	domainitem "sysken-pay-api/app/domain/object/item"
	domainpurchase "sysken-pay-api/app/domain/object/purchase"
)

// --- mocks ---

type mockPurchaseRepo struct {
	createFunc func(ctx context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error)
	cancelFunc func(ctx context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error)
}

func (m *mockPurchaseRepo) CreatePurchase(ctx context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
	return m.createFunc(ctx, p)
}

func (m *mockPurchaseRepo) CancelPurchase(ctx context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
	return m.cancelFunc(ctx, p)
}

type mockItemRepo struct {
	getByIDFunc      func(ctx context.Context, id int) (*domainitem.Item, error)
	getByJanCodeFunc func(ctx context.Context, janCode string) (*domainitem.Item, error)
	insertFunc       func(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error)
	updateFunc       func(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error)
	getAllFunc        func(ctx context.Context) ([]*domainitem.Item, error)
}

func (m *mockItemRepo) GetItemByID(ctx context.Context, id int) (*domainitem.Item, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockItemRepo) GetItemByJanCode(ctx context.Context, janCode string) (*domainitem.Item, error) {
	return m.getByJanCodeFunc(ctx, janCode)
}
func (m *mockItemRepo) InsertItem(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error) {
	return m.insertFunc(ctx, item)
}
func (m *mockItemRepo) UpdateItem(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error) {
	return m.updateFunc(ctx, item)
}
func (m *mockItemRepo) GetAllItems(ctx context.Context) ([]*domainitem.Item, error) {
	return m.getAllFunc(ctx)
}

type mockBalanceRepo struct {
	getBalanceFunc          func(ctx context.Context, userID string) (*domainbalance.Balance, error)
	getPurchaseHistoryFunc  func(ctx context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error)
}

func (m *mockBalanceRepo) GetBalance(ctx context.Context, userID string) (*domainbalance.Balance, error) {
	return m.getBalanceFunc(ctx, userID)
}
func (m *mockBalanceRepo) GetPurchaseHistories(ctx context.Context, userID string, limit, offset int) ([]*domainbalance.PurchaseHistory, int, int, error) {
	return m.getPurchaseHistoryFunc(ctx, userID, limit, offset)
}

type mockTxManager struct{}

func (m *mockTxManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

// helper: valid item for tests (price=200)
func newTestItem(price int) *domainitem.Item {
	i, _ := domainitem.NewItem("4901234567894", "テスト商品", price)
	return i
}

// helper: valid balance
func newTestBalance(userID string, balance int) *domainbalance.Balance {
	b, _ := domainbalance.NewBalance(userID, balance)
	return b
}

// --- CreatePurchase ---

func TestCreatePurchase_Success(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 2}}

	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return newTestItem(200), nil
		},
	}
	balanceRepo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return newTestBalance(userID, 1000), nil // 残高1000円 >= 合計400円
		},
	}
	purchaseRepo := &mockPurchaseRepo{
		createFunc: func(_ context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
			return p, nil
		},
	}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	result, err := uc.CreatePurchase(context.Background(), "user-1", inputs)
	if err != nil {
		t.Fatalf("CreatePurchase should succeed: %v", err)
	}
	if result.UserID() != "user-1" {
		t.Errorf("UserID() = %s, want user-1", result.UserID())
	}
}

func TestCreatePurchase_InsufficientBalance(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 5}}

	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return newTestItem(300), nil // 合計1500円
		},
	}
	balanceRepo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return newTestBalance(userID, 1000), nil // 残高1000円 < 1500円
		},
	}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CreatePurchase should fail with insufficient balance")
	}
}

func TestCreatePurchase_EmptyUserID(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return newTestItem(100), nil
		},
	}
	balanceRepo := &mockBalanceRepo{}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "", inputs); err == nil {
		t.Error("CreatePurchase with empty userID should fail")
	}
}

func TestCreatePurchase_EmptyItems(t *testing.T) {
	itemRepo := &mockItemRepo{}
	balanceRepo := &mockBalanceRepo{}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", []PurchaseItemInput{}); err == nil {
		t.Error("CreatePurchase with empty items should fail")
	}
}

func TestCreatePurchase_InvalidItemID(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 0, Quantity: 1}}
	itemRepo := &mockItemRepo{}
	balanceRepo := &mockBalanceRepo{}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CreatePurchase with itemID=0 should fail")
	}
}

func TestCreatePurchase_ItemRepoError(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return nil, errors.New("item not found")
		},
	}
	balanceRepo := &mockBalanceRepo{}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CreatePurchase should propagate item repo error")
	}
}

func TestCreatePurchase_PurchaseRepoError(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return newTestItem(100), nil
		},
	}
	balanceRepo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return newTestBalance(userID, 500), nil
		},
	}
	purchaseRepo := &mockPurchaseRepo{
		createFunc: func(_ context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
			return nil, errors.New("db error")
		},
	}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CreatePurchase should propagate purchase repo error")
	}
}

func TestCreatePurchase_BalanceExactlyEnough(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 5}} // 5 * 200 = 1000

	itemRepo := &mockItemRepo{
		getByIDFunc: func(_ context.Context, id int) (*domainitem.Item, error) {
			return newTestItem(200), nil
		},
	}
	balanceRepo := &mockBalanceRepo{
		getBalanceFunc: func(_ context.Context, userID string) (*domainbalance.Balance, error) {
			return newTestBalance(userID, 1000), nil // 残高1000円 == 合計1000円
		},
	}
	purchaseRepo := &mockPurchaseRepo{
		createFunc: func(_ context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
			return p, nil
		},
	}

	uc := NewCreatePurchaseUseCase(purchaseRepo, itemRepo, balanceRepo, &mockTxManager{})
	if _, err := uc.CreatePurchase(context.Background(), "user-1", inputs); err != nil {
		t.Errorf("CreatePurchase with exact balance should succeed: %v", err)
	}
}

// --- CancelPurchase ---

func TestCancelPurchase_Success(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	purchaseRepo := &mockPurchaseRepo{
		cancelFunc: func(_ context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
			return p, nil
		},
	}

	uc := NewCancelPurchaseUseCase(purchaseRepo)
	result, err := uc.CancelPurchase(context.Background(), "user-1", inputs)
	if err != nil {
		t.Fatalf("CancelPurchase should succeed: %v", err)
	}
	if result.UserID() != "user-1" {
		t.Errorf("UserID() = %s, want user-1", result.UserID())
	}
}

func TestCancelPurchase_EmptyUserID(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCancelPurchaseUseCase(purchaseRepo)
	if _, err := uc.CancelPurchase(context.Background(), "", inputs); err == nil {
		t.Error("CancelPurchase with empty userID should fail")
	}
}

func TestCancelPurchase_EmptyItems(t *testing.T) {
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCancelPurchaseUseCase(purchaseRepo)
	if _, err := uc.CancelPurchase(context.Background(), "user-1", []PurchaseItemInput{}); err == nil {
		t.Error("CancelPurchase with empty items should fail")
	}
}

func TestCancelPurchase_InvalidItemID(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: -1, Quantity: 1}}
	purchaseRepo := &mockPurchaseRepo{}

	uc := NewCancelPurchaseUseCase(purchaseRepo)
	if _, err := uc.CancelPurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CancelPurchase with invalid itemID should fail")
	}
}

func TestCancelPurchase_RepoError(t *testing.T) {
	inputs := []PurchaseItemInput{{ItemID: 1, Quantity: 1}}
	purchaseRepo := &mockPurchaseRepo{
		cancelFunc: func(_ context.Context, p *domainpurchase.Purchase) (*domainpurchase.Purchase, error) {
			return nil, errors.New("db error")
		},
	}

	uc := NewCancelPurchaseUseCase(purchaseRepo)
	if _, err := uc.CancelPurchase(context.Background(), "user-1", inputs); err == nil {
		t.Error("CancelPurchase should propagate repo error")
	}
}

