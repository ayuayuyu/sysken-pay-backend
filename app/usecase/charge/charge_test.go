package charge

import (
	"context"
	"errors"
	"testing"

	domaincharge "sysken-pay-api/app/domain/object/charge"
)

// --- mock ---

type mockChargeRepo struct {
	chargeAmountFunc func(ctx context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error)
	chargeCancelFunc func(ctx context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error)
}

func (m *mockChargeRepo) ChargeAmount(ctx context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
	return m.chargeAmountFunc(ctx, c)
}

func (m *mockChargeRepo) ChargeCancel(ctx context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
	return m.chargeCancelFunc(ctx, c)
}

// --- ChargeAmount ---

func TestChargeAmount_Success(t *testing.T) {
	repo := &mockChargeRepo{
		chargeAmountFunc: func(_ context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
			c.SetBalance(1500)
			return c, nil
		},
	}
	uc := NewChargeAmountUseCase(repo)
	result, err := uc.ChargeAmount(context.Background(), "user-1", 1000)
	if err != nil {
		t.Fatalf("ChargeAmount should succeed: %v", err)
	}
	if result.Balance() != 1500 {
		t.Errorf("Balance() = %d, want 1500", result.Balance())
	}
}

func TestChargeAmount_InvalidAmount(t *testing.T) {
	repo := &mockChargeRepo{}
	uc := NewChargeAmountUseCase(repo)
	if _, err := uc.ChargeAmount(context.Background(), "user-1", 300); err == nil {
		t.Error("ChargeAmount with invalid amount should fail")
	}
}

func TestChargeAmount_EmptyUserID(t *testing.T) {
	repo := &mockChargeRepo{}
	uc := NewChargeAmountUseCase(repo)
	if _, err := uc.ChargeAmount(context.Background(), "", 1000); err == nil {
		t.Error("ChargeAmount with empty userID should fail")
	}
}

func TestChargeAmount_RepoError(t *testing.T) {
	repo := &mockChargeRepo{
		chargeAmountFunc: func(_ context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewChargeAmountUseCase(repo)
	if _, err := uc.ChargeAmount(context.Background(), "user-1", 1000); err == nil {
		t.Error("ChargeAmount should propagate repo error")
	}
}

func TestChargeAmount_AllAllowedAmounts(t *testing.T) {
	for _, amount := range []int{500, 1000, 2000} {
		repo := &mockChargeRepo{
			chargeAmountFunc: func(_ context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
				return c, nil
			},
		}
		uc := NewChargeAmountUseCase(repo)
		if _, err := uc.ChargeAmount(context.Background(), "user-1", amount); err != nil {
			t.Errorf("ChargeAmount(%d) should succeed: %v", amount, err)
		}
	}
}

// --- ChargeCancel ---

func TestChargeCancel_Success(t *testing.T) {
	repo := &mockChargeRepo{
		chargeCancelFunc: func(_ context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
			c.SetBalance(500)
			return c, nil
		},
	}
	uc := NewChargeCancelUseCase(repo)
	result, err := uc.ChargeCancel(context.Background(), "user-1", 1000)
	if err != nil {
		t.Fatalf("ChargeCancel should succeed: %v", err)
	}
	if result.Balance() != 500 {
		t.Errorf("Balance() = %d, want 500", result.Balance())
	}
}

func TestChargeCancel_InvalidAmount(t *testing.T) {
	repo := &mockChargeRepo{}
	uc := NewChargeCancelUseCase(repo)
	if _, err := uc.ChargeCancel(context.Background(), "user-1", 999); err == nil {
		t.Error("ChargeCancel with invalid amount should fail")
	}
}

func TestChargeCancel_EmptyUserID(t *testing.T) {
	repo := &mockChargeRepo{}
	uc := NewChargeCancelUseCase(repo)
	if _, err := uc.ChargeCancel(context.Background(), "", 1000); err == nil {
		t.Error("ChargeCancel with empty userID should fail")
	}
}

func TestChargeCancel_RepoError(t *testing.T) {
	repo := &mockChargeRepo{
		chargeCancelFunc: func(_ context.Context, c *domaincharge.Charge) (*domaincharge.Charge, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewChargeCancelUseCase(repo)
	if _, err := uc.ChargeCancel(context.Background(), "user-1", 1000); err == nil {
		t.Error("ChargeCancel should propagate repo error")
	}
}
