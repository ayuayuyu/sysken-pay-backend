package charge

import (
	"testing"
	"time"
)

// --- SetAmount ---

func TestSetAmount_AllowedValues(t *testing.T) {
	allowed := []int{500, 1000, 2000}
	for _, amount := range allowed {
		c := &Charge{}
		if err := c.SetAmount(amount); err != nil {
			t.Errorf("SetAmount(%d) should succeed, got: %v", amount, err)
		}
		if c.Amount() != amount {
			t.Errorf("Amount() = %d, want %d", c.Amount(), amount)
		}
	}
}

func TestSetAmount_DisallowedValues(t *testing.T) {
	disallowed := []int{0, -1, 100, 499, 501, 1001, 1500, 1999, 2001, 3000, 5000}
	for _, amount := range disallowed {
		c := &Charge{}
		if err := c.SetAmount(amount); err == nil {
			t.Errorf("SetAmount(%d) should fail", amount)
		}
	}
}

// --- SetUserID ---

func TestSetUserID_Valid(t *testing.T) {
	c := &Charge{}
	if err := c.SetUserID("user-123"); err != nil {
		t.Errorf("SetUserID should succeed: %v", err)
	}
}

func TestSetUserID_Empty(t *testing.T) {
	c := &Charge{}
	if err := c.SetUserID(""); err == nil {
		t.Error("SetUserID empty should fail")
	}
}

// --- SetBalance ---

func TestSetBalance_Zero(t *testing.T) {
	c := &Charge{}
	if err := c.SetBalance(0); err != nil {
		t.Errorf("SetBalance(0) should succeed: %v", err)
	}
}

func TestSetBalance_Positive(t *testing.T) {
	c := &Charge{}
	if err := c.SetBalance(1000); err != nil {
		t.Errorf("SetBalance(1000) should succeed: %v", err)
	}
}

func TestSetBalance_Negative(t *testing.T) {
	c := &Charge{}
	if err := c.SetBalance(-1); err == nil {
		t.Error("SetBalance(-1) should fail")
	}
}

// --- SetCreatedAt ---

func TestSetCreatedAt_Past(t *testing.T) {
	c := &Charge{}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	past := time.Now().In(jst).Add(-time.Hour)
	if err := c.SetCreatedAt(past); err != nil {
		t.Errorf("SetCreatedAt past JST should succeed: %v", err)
	}
}

func TestSetCreatedAt_Future(t *testing.T) {
	c := &Charge{}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	future := time.Now().In(jst).Add(time.Hour)
	if err := c.SetCreatedAt(future); err == nil {
		t.Error("SetCreatedAt future should fail")
	}
}

// --- NewCharge ---

func TestNewCharge_Valid(t *testing.T) {
	tests := []struct {
		userID string
		amount int
	}{
		{"user-1", 500},
		{"user-2", 1000},
		{"user-3", 2000},
	}
	for _, tt := range tests {
		c, err := NewCharge(tt.userID, tt.amount)
		if err != nil {
			t.Errorf("NewCharge(%s,%d) should succeed: %v", tt.userID, tt.amount, err)
			continue
		}
		if c.UserID() != tt.userID {
			t.Errorf("UserID() = %s, want %s", c.UserID(), tt.userID)
		}
		if c.Amount() != tt.amount {
			t.Errorf("Amount() = %d, want %d", c.Amount(), tt.amount)
		}
	}
}

func TestNewCharge_InvalidAmount(t *testing.T) {
	if _, err := NewCharge("user-1", 300); err == nil {
		t.Error("NewCharge with invalid amount should fail")
	}
}

func TestNewCharge_EmptyUserID(t *testing.T) {
	if _, err := NewCharge("", 1000); err == nil {
		t.Error("NewCharge with empty userID should fail")
	}
}
