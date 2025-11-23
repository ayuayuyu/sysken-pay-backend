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
	updatedAt time.Time
	deletedAt time.Time
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

func (u *User) SetCreatedAt(createdAt time.Time) error {

	// createdAtは未来の日付でないこと
	if createdAt.After(time.Now()) {
		return errors.New("createdAt must not be in the future")
	}

	// createdAtのタイムゾーンはJSTであること
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstCreatedAt := createdAt.In(jst)
	if !createdAt.Equal(jstCreatedAt) {
		return errors.New("createdAt must be in JST")
	}

	u.createdAt = createdAt

	return nil
}

func (u *User) SetUpdatedAt(updatedAt time.Time) error {

	// updatedAtは未来の日付でないこと
	if updatedAt.After(time.Now()) {
		return errors.New("updatedAt must not be in the future")
	}

	// updatedAtのタイムゾーンはJSTであること
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstUpdatedAt := updatedAt.In(jst)
	if !updatedAt.Equal(jstUpdatedAt) {
		return errors.New("updatedAt must be in JST")
	}

	u.updatedAt = updatedAt

	return nil
}

func (u *User) SetDeletedAt(deletedAt time.Time) error {

	// deletedAtは未来の日付でないこと
	if deletedAt.After(time.Now()) {
		return errors.New("deletedAt must not be in the future")
	}

	// deletedAtのタイムゾーンはJSTであること
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstDeletedAt := deletedAt.In(jst)
	if !deletedAt.Equal(jstDeletedAt) {
		return errors.New("deletedAt must be in JST")
	}

	u.deletedAt = deletedAt

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
func (p *User) UpdatedAt() time.Time {
	return p.updatedAt
}
func (p *User) DeletedAt() time.Time {
	return p.deletedAt
}
