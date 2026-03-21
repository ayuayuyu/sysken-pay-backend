package user

import (
	"strings"
	"testing"
)

// --- SetUserID ---

func TestSetUserID_Valid(t *testing.T) {
	u := &User{}
	if err := u.SetUserID("user123"); err != nil {
		t.Errorf("SetUserID should succeed: %v", err)
	}
}

func TestSetUserID_Empty(t *testing.T) {
	u := &User{}
	if err := u.SetUserID(""); err == nil {
		t.Error("SetUserID empty should fail")
	}
}

func TestSetUserID_Exactly20Chars(t *testing.T) {
	u := &User{}
	id := strings.Repeat("a", 20)
	if err := u.SetUserID(id); err != nil {
		t.Errorf("SetUserID 20 chars should succeed: %v", err)
	}
}

func TestSetUserID_Over20Chars(t *testing.T) {
	u := &User{}
	id := strings.Repeat("a", 21)
	if err := u.SetUserID(id); err == nil {
		t.Error("SetUserID 21 chars should fail")
	}
}

func TestSetUserID_MultibyteCounts(t *testing.T) {
	u := &User{}
	// 日本語20文字はOK
	id := strings.Repeat("あ", 20)
	if err := u.SetUserID(id); err != nil {
		t.Errorf("SetUserID 20 multibyte chars should succeed: %v", err)
	}
	// 21文字はNG
	id21 := strings.Repeat("あ", 21)
	if err := u.SetUserID(id21); err == nil {
		t.Error("SetUserID 21 multibyte chars should fail")
	}
}

// --- SetUserName ---

func TestSetUserName_Valid(t *testing.T) {
	u := &User{}
	if err := u.SetUserName("田中 太郎"); err != nil {
		t.Errorf("SetUserName should succeed: %v", err)
	}
}

func TestSetUserName_Empty(t *testing.T) {
	u := &User{}
	if err := u.SetUserName(""); err == nil {
		t.Error("SetUserName empty should fail")
	}
}

func TestSetUserName_Exactly50Chars(t *testing.T) {
	u := &User{}
	name := strings.Repeat("a", 50)
	if err := u.SetUserName(name); err != nil {
		t.Errorf("SetUserName 50 chars should succeed: %v", err)
	}
}

func TestSetUserName_Over50Chars(t *testing.T) {
	u := &User{}
	name := strings.Repeat("a", 51)
	if err := u.SetUserName(name); err == nil {
		t.Error("SetUserName 51 chars should fail")
	}
}

func TestSetUserName_MultibyteCounts(t *testing.T) {
	u := &User{}
	name50 := strings.Repeat("あ", 50)
	if err := u.SetUserName(name50); err != nil {
		t.Errorf("SetUserName 50 multibyte chars should succeed: %v", err)
	}
	name51 := strings.Repeat("あ", 51)
	if err := u.SetUserName(name51); err == nil {
		t.Error("SetUserName 51 multibyte chars should fail")
	}
}

// --- NewUser ---

func TestNewUser_Valid(t *testing.T) {
	u, err := NewUser("student001", "田中 太郎")
	if err != nil {
		t.Fatalf("NewUser should succeed: %v", err)
	}
	if u.ID() != "student001" {
		t.Errorf("ID() = %s, want student001", u.ID())
	}
	if u.UserName() != "田中 太郎" {
		t.Errorf("UserName() = %s, want 田中 太郎", u.UserName())
	}
}

func TestNewUser_EmptyUserID(t *testing.T) {
	if _, err := NewUser("", "田中 太郎"); err == nil {
		t.Error("NewUser empty userID should fail")
	}
}

func TestNewUser_EmptyUserName(t *testing.T) {
	if _, err := NewUser("student001", ""); err == nil {
		t.Error("NewUser empty userName should fail")
	}
}

func TestNewUser_UserIDTooLong(t *testing.T) {
	if _, err := NewUser(strings.Repeat("a", 21), "田中 太郎"); err == nil {
		t.Error("NewUser userID > 20 chars should fail")
	}
}

func TestNewUser_UserNameTooLong(t *testing.T) {
	if _, err := NewUser("student001", strings.Repeat("a", 51)); err == nil {
		t.Error("NewUser userName > 50 chars should fail")
	}
}
