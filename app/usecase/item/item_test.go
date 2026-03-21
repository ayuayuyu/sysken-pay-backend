package item

import (
	"context"
	"errors"
	"testing"

	domainitem "sysken-pay-api/app/domain/object/item"
)

// --- mock ---

type mockItemRepo struct {
	insertFunc       func(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error)
	updateFunc       func(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error)
	getByIDFunc      func(ctx context.Context, id int) (*domainitem.Item, error)
	getByJanCodeFunc func(ctx context.Context, janCode string) (*domainitem.Item, error)
	getAllFunc        func(ctx context.Context) ([]*domainitem.Item, error)
}

func (m *mockItemRepo) InsertItem(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error) {
	return m.insertFunc(ctx, item)
}
func (m *mockItemRepo) UpdateItem(ctx context.Context, item *domainitem.Item) (*domainitem.Item, error) {
	return m.updateFunc(ctx, item)
}
func (m *mockItemRepo) GetItemByID(ctx context.Context, id int) (*domainitem.Item, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockItemRepo) GetItemByJanCode(ctx context.Context, janCode string) (*domainitem.Item, error) {
	return m.getByJanCodeFunc(ctx, janCode)
}
func (m *mockItemRepo) GetAllItems(ctx context.Context) ([]*domainitem.Item, error) {
	return m.getAllFunc(ctx)
}

const validJAN13 = "4901234567894"
const validJAN8 = "49012347"

// --- RegisterItem ---

func TestRegisterItem_Success(t *testing.T) {
	repo := &mockItemRepo{
		insertFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return item, nil
		},
	}
	uc := NewRegisterItemUseCase(repo)
	result, err := uc.RegisterItem(context.Background(), validJAN13, "コーラ", 150)
	if err != nil {
		t.Fatalf("RegisterItem should succeed: %v", err)
	}
	if result.JanCode() != validJAN13 {
		t.Errorf("JanCode() = %s, want %s", result.JanCode(), validJAN13)
	}
}

func TestRegisterItem_InvalidJanCode(t *testing.T) {
	repo := &mockItemRepo{}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), "1234567890123", "コーラ", 150); err == nil {
		t.Error("RegisterItem with invalid JAN code should fail")
	}
}

func TestRegisterItem_EmptyItemName(t *testing.T) {
	repo := &mockItemRepo{}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), validJAN13, "", 150); err == nil {
		t.Error("RegisterItem with empty item name should fail")
	}
}

func TestRegisterItem_NegativePrice(t *testing.T) {
	repo := &mockItemRepo{}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), validJAN13, "コーラ", -1); err == nil {
		t.Error("RegisterItem with negative price should fail")
	}
}

func TestRegisterItem_ZeroPrice(t *testing.T) {
	repo := &mockItemRepo{
		insertFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return item, nil
		},
	}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), validJAN13, "無料サンプル", 0); err != nil {
		t.Errorf("RegisterItem with price=0 should succeed: %v", err)
	}
}

func TestRegisterItem_RepoError(t *testing.T) {
	repo := &mockItemRepo{
		insertFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), validJAN13, "コーラ", 150); err == nil {
		t.Error("RegisterItem should propagate repo error")
	}
}

func TestRegisterItem_8DigitJAN(t *testing.T) {
	repo := &mockItemRepo{
		insertFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return item, nil
		},
	}
	uc := NewRegisterItemUseCase(repo)
	if _, err := uc.RegisterItem(context.Background(), validJAN8, "商品", 100); err != nil {
		t.Errorf("RegisterItem with 8-digit JAN should succeed: %v", err)
	}
}

// --- UpdateItem ---

func TestUpdateItem_Success(t *testing.T) {
	repo := &mockItemRepo{
		updateFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return item, nil
		},
	}
	uc := NewUpdateItemUseCase(repo)
	result, err := uc.UpdateItem(context.Background(), validJAN13, "新コーラ", 200)
	if err != nil {
		t.Fatalf("UpdateItem should succeed: %v", err)
	}
	if result.JanCode() != validJAN13 {
		t.Errorf("JanCode() = %s, want %s", result.JanCode(), validJAN13)
	}
}

func TestUpdateItem_InvalidJanCode(t *testing.T) {
	repo := &mockItemRepo{}
	uc := NewUpdateItemUseCase(repo)
	if _, err := uc.UpdateItem(context.Background(), "badcode", "コーラ", 150); err == nil {
		t.Error("UpdateItem with invalid JAN code should fail")
	}
}

func TestUpdateItem_EmptyItemName(t *testing.T) {
	repo := &mockItemRepo{}
	uc := NewUpdateItemUseCase(repo)
	if _, err := uc.UpdateItem(context.Background(), validJAN13, "", 150); err == nil {
		t.Error("UpdateItem with empty item name should fail")
	}
}

func TestUpdateItem_RepoError(t *testing.T) {
	repo := &mockItemRepo{
		updateFunc: func(_ context.Context, item *domainitem.Item) (*domainitem.Item, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewUpdateItemUseCase(repo)
	if _, err := uc.UpdateItem(context.Background(), validJAN13, "コーラ", 150); err == nil {
		t.Error("UpdateItem should propagate repo error")
	}
}

// --- FindItemByJanCode ---

func TestFindItemByJanCode_Success(t *testing.T) {
	repo := &mockItemRepo{
		getByJanCodeFunc: func(_ context.Context, janCode string) (*domainitem.Item, error) {
			i, _ := domainitem.NewItem(janCode, "コーラ", 150)
			return i, nil
		},
	}
	uc := NewFindItemByJanCodeUseCase(repo)
	result, err := uc.GetItemByJanCode(context.Background(), validJAN13)
	if err != nil {
		t.Fatalf("GetItemByJanCode should succeed: %v", err)
	}
	if result.JanCode() != validJAN13 {
		t.Errorf("JanCode() = %s, want %s", result.JanCode(), validJAN13)
	}
}

func TestFindItemByJanCode_RepoError(t *testing.T) {
	repo := &mockItemRepo{
		getByJanCodeFunc: func(_ context.Context, janCode string) (*domainitem.Item, error) {
			return nil, errors.New("not found")
		},
	}
	uc := NewFindItemByJanCodeUseCase(repo)
	if _, err := uc.GetItemByJanCode(context.Background(), validJAN13); err == nil {
		t.Error("GetItemByJanCode should propagate repo error")
	}
}

// --- GetAllItems ---

func TestGetAllItems_Success(t *testing.T) {
	i1, _ := domainitem.NewItem(validJAN13, "コーラ", 150)
	i2, _ := domainitem.NewItem(validJAN8, "パン", 200)
	repo := &mockItemRepo{
		getAllFunc: func(_ context.Context) ([]*domainitem.Item, error) {
			return []*domainitem.Item{i1, i2}, nil
		},
	}
	uc := NewGetAllItemsUseCase(repo)
	result, err := uc.GetAllItems(context.Background())
	if err != nil {
		t.Fatalf("GetAllItems should succeed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len(result) = %d, want 2", len(result))
	}
}

func TestGetAllItems_Empty(t *testing.T) {
	repo := &mockItemRepo{
		getAllFunc: func(_ context.Context) ([]*domainitem.Item, error) {
			return []*domainitem.Item{}, nil
		},
	}
	uc := NewGetAllItemsUseCase(repo)
	result, err := uc.GetAllItems(context.Background())
	if err != nil {
		t.Fatalf("GetAllItems should succeed with empty list: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("len(result) = %d, want 0", len(result))
	}
}

func TestGetAllItems_RepoError(t *testing.T) {
	repo := &mockItemRepo{
		getAllFunc: func(_ context.Context) ([]*domainitem.Item, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewGetAllItemsUseCase(repo)
	if _, err := uc.GetAllItems(context.Background()); err == nil {
		t.Error("GetAllItems should propagate repo error")
	}
}
