package item

import (
	"strings"
	"testing"
	"time"
)

// 有効なJANコード（テスト用に検証済み）
// 13桁: "4901234567894" (チェックディジット=4)
// 8桁:  "49012347"     (チェックディジット=7)
const (
	validJAN13 = "4901234567894"
	validJAN8  = "49012347"
)

// --- SetJanCode ---

func TestSetJanCode_Valid13(t *testing.T) {
	p := &Item{}
	if err := p.SetJanCode(validJAN13); err != nil {
		t.Errorf("SetJanCode valid 13-digit should succeed: %v", err)
	}
}

func TestSetJanCode_Valid8(t *testing.T) {
	p := &Item{}
	if err := p.SetJanCode(validJAN8); err != nil {
		t.Errorf("SetJanCode valid 8-digit should succeed: %v", err)
	}
}

func TestSetJanCode_Empty(t *testing.T) {
	p := &Item{}
	if err := p.SetJanCode(""); err == nil {
		t.Error("SetJanCode empty should fail")
	}
}

func TestSetJanCode_WrongLength(t *testing.T) {
	for _, code := range []string{"490123", "4901234567", "49012345678901"} {
		p := &Item{}
		if err := p.SetJanCode(code); err == nil {
			t.Errorf("SetJanCode length=%d should fail", len(code))
		}
	}
}

func TestSetJanCode_NonNumeric(t *testing.T) {
	p := &Item{}
	if err := p.SetJanCode("490123456789A"); err == nil {
		t.Error("SetJanCode non-numeric should fail")
	}
}

func TestSetJanCode_InvalidCheckDigit(t *testing.T) {
	p := &Item{}
	// 最後の桁を変えて無効にする
	invalid := validJAN13[:12] + "0"
	if err := p.SetJanCode(invalid); err == nil {
		t.Error("SetJanCode invalid check digit should fail")
	}
}

// --- SetName ---

func TestSetName_Valid(t *testing.T) {
	p := &Item{}
	if err := p.SetName("コーラ"); err != nil {
		t.Errorf("SetName should succeed: %v", err)
	}
}

func TestSetName_Empty(t *testing.T) {
	p := &Item{}
	if err := p.SetName(""); err == nil {
		t.Error("SetName empty should fail")
	}
}

// --- SetPrice ---

func TestSetPrice_Zero(t *testing.T) {
	p := &Item{}
	if err := p.SetPrice(0); err != nil {
		t.Errorf("SetPrice(0) should succeed: %v", err)
	}
}

func TestSetPrice_Positive(t *testing.T) {
	p := &Item{}
	if err := p.SetPrice(150); err != nil {
		t.Errorf("SetPrice(150) should succeed: %v", err)
	}
}

func TestSetPrice_Negative(t *testing.T) {
	p := &Item{}
	if err := p.SetPrice(-1); err == nil {
		t.Error("SetPrice(-1) should fail")
	}
}

// --- SetID ---

func TestSetID_Valid(t *testing.T) {
	p := &Item{}
	if err := p.SetID(1); err != nil {
		t.Errorf("SetID(1) should succeed: %v", err)
	}
}

func TestSetID_Zero(t *testing.T) {
	p := &Item{}
	if err := p.SetID(0); err == nil {
		t.Error("SetID(0) should fail")
	}
}

func TestSetID_Negative(t *testing.T) {
	p := &Item{}
	if err := p.SetID(-1); err == nil {
		t.Error("SetID(-1) should fail")
	}
}

// --- SetCreatedAt / SetUpdatedAt ---

func TestSetCreatedAt_PastJST(t *testing.T) {
	p := &Item{}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	if err := p.SetCreatedAt(time.Now().In(jst).Add(-time.Hour)); err != nil {
		t.Errorf("SetCreatedAt past JST should succeed: %v", err)
	}
}

func TestSetCreatedAt_Future(t *testing.T) {
	p := &Item{}
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	if err := p.SetCreatedAt(time.Now().In(jst).Add(time.Hour)); err == nil {
		t.Error("SetCreatedAt future should fail")
	}
}

// --- NewItem ---

func TestNewItem_Valid(t *testing.T) {
	i, err := NewItem(validJAN13, "コーラ", 150)
	if err != nil {
		t.Fatalf("NewItem should succeed: %v", err)
	}
	if i.JanCode() != validJAN13 {
		t.Errorf("JanCode() = %s, want %s", i.JanCode(), validJAN13)
	}
	if i.Name() != "コーラ" {
		t.Errorf("Name() = %s, want コーラ", i.Name())
	}
	if i.Price() != 150 {
		t.Errorf("Price() = %d, want 150", i.Price())
	}
}

func TestNewItem_InvalidJanCode(t *testing.T) {
	if _, err := NewItem("invalid", "コーラ", 150); err == nil {
		t.Error("NewItem with invalid JAN code should fail")
	}
}

func TestNewItem_EmptyName(t *testing.T) {
	if _, err := NewItem(validJAN13, "", 150); err == nil {
		t.Error("NewItem with empty name should fail")
	}
}

func TestNewItem_NegativePrice(t *testing.T) {
	if _, err := NewItem(validJAN13, "コーラ", -1); err == nil {
		t.Error("NewItem with negative price should fail")
	}
}

// --- isValidJanCode (間接検証) ---

func TestIsValidJanCode_AllDigitsSame(t *testing.T) {
	// "0000000" は8桁でチェックディジット0
	// sum = 0 → check = (10-0)%10 = 0
	p := &Item{}
	if err := p.SetJanCode("00000000"); err != nil {
		t.Errorf("00000000 should be valid: %v", err)
	}
}

// --- NewItem name boundary ---

func TestNewItem_LongName(t *testing.T) {
	longName := strings.Repeat("あ", 100)
	i, err := NewItem(validJAN13, longName, 100)
	if err != nil {
		t.Errorf("NewItem long name should succeed: %v", err)
	}
	if i.Name() != longName {
		t.Error("Name mismatch")
	}
}
