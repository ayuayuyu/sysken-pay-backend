package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainuser "sysken-pay-api/app/domain/object/user"
)

// --- mock ---

type mockUserRepo struct {
	insertFunc func(ctx context.Context, u *domainuser.User) (*domainuser.User, error)
	updateFunc func(ctx context.Context, u *domainuser.User) (*domainuser.User, error)
}

func (m *mockUserRepo) InsertUser(ctx context.Context, u *domainuser.User) (*domainuser.User, error) {
	return m.insertFunc(ctx, u)
}

func (m *mockUserRepo) UpdateUser(ctx context.Context, u *domainuser.User) (*domainuser.User, error) {
	return m.updateFunc(ctx, u)
}

// --- RegisterUser ---

func TestRegisterUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		insertFunc: func(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
			return u, nil
		},
	}
	uc := NewRegisterUserUseCase(repo)
	result, err := uc.RegisterUser(context.Background(), "student001", "田中 太郎")
	if err != nil {
		t.Fatalf("RegisterUser should succeed: %v", err)
	}
	if result.ID() != "student001" {
		t.Errorf("ID() = %s, want student001", result.ID())
	}
	if result.UserName() != "田中 太郎" {
		t.Errorf("UserName() = %s, want 田中 太郎", result.UserName())
	}
}

func TestRegisterUser_EmptyUserID(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), "", "田中 太郎"); err == nil {
		t.Error("RegisterUser with empty userID should fail")
	}
}

func TestRegisterUser_EmptyUserName(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), "student001", ""); err == nil {
		t.Error("RegisterUser with empty userName should fail")
	}
}

func TestRegisterUser_UserIDTooLong(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), strings.Repeat("a", 21), "田中 太郎"); err == nil {
		t.Error("RegisterUser with userID > 20 chars should fail")
	}
}

func TestRegisterUser_UserNameTooLong(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), "student001", strings.Repeat("a", 51)); err == nil {
		t.Error("RegisterUser with userName > 50 chars should fail")
	}
}

func TestRegisterUser_RepoError(t *testing.T) {
	repo := &mockUserRepo{
		insertFunc: func(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), "student001", "田中 太郎"); err == nil {
		t.Error("RegisterUser should propagate repo error")
	}
}

func TestRegisterUser_MaxLengthUserID(t *testing.T) {
	repo := &mockUserRepo{
		insertFunc: func(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
			return u, nil
		},
	}
	uc := NewRegisterUserUseCase(repo)
	if _, err := uc.RegisterUser(context.Background(), strings.Repeat("a", 20), "田中 太郎"); err != nil {
		t.Errorf("RegisterUser with 20-char userID should succeed: %v", err)
	}
}

// --- UpdateUser ---

func TestUpdateUser_Success(t *testing.T) {
	repo := &mockUserRepo{
		updateFunc: func(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
			return u, nil
		},
	}
	uc := NewUpdateUserUseCase(repo)
	result, err := uc.UpdateUser(context.Background(), "student001", "佐藤 花子")
	if err != nil {
		t.Fatalf("UpdateUser should succeed: %v", err)
	}
	if result.ID() != "student001" {
		t.Errorf("ID() = %s, want student001", result.ID())
	}
	if result.UserName() != "佐藤 花子" {
		t.Errorf("UserName() = %s, want 佐藤 花子", result.UserName())
	}
}

func TestUpdateUser_EmptyUserName(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewUpdateUserUseCase(repo)
	if _, err := uc.UpdateUser(context.Background(), "student001", ""); err == nil {
		t.Error("UpdateUser with empty userName should fail")
	}
}

func TestUpdateUser_EmptyUserID(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewUpdateUserUseCase(repo)
	if _, err := uc.UpdateUser(context.Background(), "", "田中 太郎"); err == nil {
		t.Error("UpdateUser with empty userID should fail")
	}
}

func TestUpdateUser_UserNameTooLong(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewUpdateUserUseCase(repo)
	if _, err := uc.UpdateUser(context.Background(), "student001", strings.Repeat("a", 51)); err == nil {
		t.Error("UpdateUser with userName > 50 chars should fail")
	}
}

func TestUpdateUser_RepoError(t *testing.T) {
	repo := &mockUserRepo{
		updateFunc: func(_ context.Context, u *domainuser.User) (*domainuser.User, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewUpdateUserUseCase(repo)
	if _, err := uc.UpdateUser(context.Background(), "student001", "田中 太郎"); err == nil {
		t.Error("UpdateUser should propagate repo error")
	}
}
