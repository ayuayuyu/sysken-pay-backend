package balance

import (
	"testing"
	"time"
)

// --- NewBalance ---

func TestNewBalance_Valid(t *testing.T) {
	b, err := NewBalance("user-1", 1500)
	if err != nil {
		t.Fatalf("NewBalance should succeed: %v", err)
	}
	if b.UserID() != "user-1" {
		t.Errorf("UserID() = %s, want user-1", b.UserID())
	}
	if b.Balance() != 1500 {
		t.Errorf("Balance() = %d, want 1500", b.Balance())
	}
}

func TestNewBalance_Zero(t *testing.T) {
	b, err := NewBalance("user-1", 0)
	if err != nil {
		t.Errorf("NewBalance(0) should succeed: %v", err)
	}
	if b.Balance() != 0 {
		t.Errorf("Balance() = %d, want 0", b.Balance())
	}
}

func TestNewBalance_EmptyUserID(t *testing.T) {
	if _, err := NewBalance("", 1000); err == nil {
		t.Error("NewBalance empty userID should fail")
	}
}

func TestNewBalance_Negative(t *testing.T) {
	if _, err := NewBalance("user-1", -1); err == nil {
		t.Error("NewBalance negative balance should fail")
	}
}

// --- NewPurchaseHistory ---

func TestNewPurchaseHistory_Valid(t *testing.T) {
	now := time.Now()
	h, err := NewPurchaseHistory("コーラ", 2, 150, now)
	if err != nil {
		t.Fatalf("NewPurchaseHistory should succeed: %v", err)
	}
	if h.ItemName() != "コーラ" {
		t.Errorf("ItemName() = %s, want コーラ", h.ItemName())
	}
	if h.Quantity() != 2 {
		t.Errorf("Quantity() = %d, want 2", h.Quantity())
	}
	if h.Price() != 150 {
		t.Errorf("Price() = %d, want 150", h.Price())
	}
	if !h.PurchaseAt().Equal(now) {
		t.Error("PurchaseAt mismatch")
	}
}

func TestNewPurchaseHistory_EmptyItemName(t *testing.T) {
	if _, err := NewPurchaseHistory("", 1, 150, time.Now()); err == nil {
		t.Error("NewPurchaseHistory empty itemName should fail")
	}
}

func TestNewPurchaseHistory_ZeroQuantity(t *testing.T) {
	if _, err := NewPurchaseHistory("コーラ", 0, 150, time.Now()); err == nil {
		t.Error("NewPurchaseHistory quantity=0 should fail")
	}
}

func TestNewPurchaseHistory_NegativeQuantity(t *testing.T) {
	if _, err := NewPurchaseHistory("コーラ", -1, 150, time.Now()); err == nil {
		t.Error("NewPurchaseHistory quantity<0 should fail")
	}
}

func TestNewPurchaseHistory_ZeroPrice(t *testing.T) {
	// 0円の商品はありうる（無料配布など）
	h, err := NewPurchaseHistory("サンプル", 1, 0, time.Now())
	if err != nil {
		t.Errorf("NewPurchaseHistory price=0 should succeed: %v", err)
	}
	if h.Price() != 0 {
		t.Errorf("Price() = %d, want 0", h.Price())
	}
}

func TestNewPurchaseHistory_NegativePrice(t *testing.T) {
	if _, err := NewPurchaseHistory("コーラ", 1, -1, time.Now()); err == nil {
		t.Error("NewPurchaseHistory negative price should fail")
	}
}
