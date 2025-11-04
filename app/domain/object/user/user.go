package user

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

//TODO モデル（データベースに入れる型を宣言する）
//データベースの制約通りになるようにエラーハンドリングをガチる
//ユーザーID、名前、作成日時、更新日時など

type User struct {
	userID    uuid.UUID
	userName  string
	createdAt time.Time
	deletedAt time.Time
}

func NewUser(userID uuid.UUID, userName string) (*User, error) {
	user := &User{}

	if err := user.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := user.SetUserName(userName); err != nil {
		return nil, err
	}
	user.createdAt = time.Now()

	return user, nil
}

func (p *User) SetUserID(userID uuid.UUID) error {
	// userIDは空でないこと
	if userID == uuid.Nil {
		return errors.New("userID must not be empty")
	}

	p.userID = userID
	return nil
}

func (p *User) SetUserName(userName string) error {
	// 表示名は1文字以上であること
	if !(utf8.RuneCountInString(userName) >= 1) {
		return errors.New("display name must be more than 0 characters")
	}

	// 表示名は50文字以下であること
	if !(utf8.RuneCountInString(userName) <= 50) {
		return errors.New("display name must be less than or equal to 50 characters")
	}

	p.userName = userName
	return nil
}

func (p *User) ID() uuid.UUID {
	return p.userID
}

func (p *User) UserName() string {
	return p.userName
}
func (p *User) CreatedAt() time.Time {
	return p.createdAt
}
func (p *User) DeletedAt() time.Time {
	return p.deletedAt
}
