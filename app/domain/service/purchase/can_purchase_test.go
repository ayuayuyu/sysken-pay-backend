package purchase

import "testing"

func TestCanPurchase_SufficientBalance(t *testing.T) {
	if err := CanPurchase(2000, 1500); err != nil {
		t.Errorf("CanPurchase should succeed when balance > total: %v", err)
	}
}

func TestCanPurchase_ExactBalance(t *testing.T) {
	if err := CanPurchase(1500, 1500); err != nil {
		t.Errorf("CanPurchase should succeed when balance == total: %v", err)
	}
}

func TestCanPurchase_InsufficientBalance(t *testing.T) {
	if err := CanPurchase(1000, 1500); err == nil {
		t.Error("CanPurchase should fail when balance < total")
	}
}

func TestCanPurchase_ZeroTotal(t *testing.T) {
	if err := CanPurchase(0, 0); err != nil {
		t.Errorf("CanPurchase(0,0) should succeed: %v", err)
	}
}

func TestCanPurchase_ZeroBalance(t *testing.T) {
	if err := CanPurchase(0, 1); err == nil {
		t.Error("CanPurchase should fail when balance=0 and total>0")
	}
}
