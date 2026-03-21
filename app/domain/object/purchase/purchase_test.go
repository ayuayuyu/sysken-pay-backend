package purchase

import "testing"

// --- NewPurchaseItem ---

func TestNewPurchaseItem_Valid(t *testing.T) {
	pi, err := NewPurchaseItem(1, 3)
	if err != nil {
		t.Fatalf("NewPurchaseItem should succeed: %v", err)
	}
	if pi.ItemID() != 1 {
		t.Errorf("ItemID() = %d, want 1", pi.ItemID())
	}
	if pi.Quantity() != 3 {
		t.Errorf("Quantity() = %d, want 3", pi.Quantity())
	}
}

func TestNewPurchaseItem_ZeroItemID(t *testing.T) {
	if _, err := NewPurchaseItem(0, 1); err == nil {
		t.Error("NewPurchaseItem itemID=0 should fail")
	}
}

func TestNewPurchaseItem_NegativeItemID(t *testing.T) {
	if _, err := NewPurchaseItem(-1, 1); err == nil {
		t.Error("NewPurchaseItem itemID<0 should fail")
	}
}

func TestNewPurchaseItem_ZeroQuantity(t *testing.T) {
	if _, err := NewPurchaseItem(1, 0); err == nil {
		t.Error("NewPurchaseItem quantity=0 should fail")
	}
}

func TestNewPurchaseItem_NegativeQuantity(t *testing.T) {
	if _, err := NewPurchaseItem(1, -1); err == nil {
		t.Error("NewPurchaseItem quantity<0 should fail")
	}
}

// --- NewPurchase ---

func TestNewPurchase_Valid(t *testing.T) {
	items := []PurchaseItem{
		{itemID: 1, quantity: 2},
		{itemID: 2, quantity: 1},
	}
	p, err := NewPurchase("user-1", items)
	if err != nil {
		t.Fatalf("NewPurchase should succeed: %v", err)
	}
	if p.UserID() != "user-1" {
		t.Errorf("UserID() = %s, want user-1", p.UserID())
	}
	if len(p.Items()) != 2 {
		t.Errorf("Items() len = %d, want 2", len(p.Items()))
	}
}

func TestNewPurchase_EmptyUserID(t *testing.T) {
	items := []PurchaseItem{{itemID: 1, quantity: 1}}
	if _, err := NewPurchase("", items); err == nil {
		t.Error("NewPurchase empty userID should fail")
	}
}

func TestNewPurchase_EmptyItems(t *testing.T) {
	if _, err := NewPurchase("user-1", []PurchaseItem{}); err == nil {
		t.Error("NewPurchase empty items should fail")
	}
}

// --- DeletePurchase ---

func TestDeletePurchase_Valid(t *testing.T) {
	items := []PurchaseItem{{itemID: 1, quantity: 1}}
	p, err := DeletePurchase("user-1", items)
	if err != nil {
		t.Fatalf("DeletePurchase should succeed: %v", err)
	}
	if p.UserID() != "user-1" {
		t.Errorf("UserID() = %s, want user-1", p.UserID())
	}
}

func TestDeletePurchase_EmptyUserID(t *testing.T) {
	items := []PurchaseItem{{itemID: 1, quantity: 1}}
	if _, err := DeletePurchase("", items); err == nil {
		t.Error("DeletePurchase empty userID should fail")
	}
}

func TestDeletePurchase_EmptyItems(t *testing.T) {
	if _, err := DeletePurchase("user-1", []PurchaseItem{}); err == nil {
		t.Error("DeletePurchase empty items should fail")
	}
}

// --- SetBalance ---

func TestSetBalance_Zero(t *testing.T) {
	p := &Purchase{}
	if err := p.SetBalance(0); err != nil {
		t.Errorf("SetBalance(0) should succeed: %v", err)
	}
}

func TestSetBalance_Positive(t *testing.T) {
	p := &Purchase{}
	if err := p.SetBalance(1000); err != nil {
		t.Errorf("SetBalance(1000) should succeed: %v", err)
	}
}

func TestSetBalance_Negative(t *testing.T) {
	p := &Purchase{}
	if err := p.SetBalance(-1); err == nil {
		t.Error("SetBalance(-1) should fail")
	}
}
