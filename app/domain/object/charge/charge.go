package charge

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//TODO モデル（データベースに入れる型を宣言する）
//データベースの制約通りになるようにエラーハンドリングをガチる

type Charge struct {
	id        int
	userID    uuid.UUID
	amount    int
	balance   int
	createdAt time.Time
	deletedAt time.Time
}

func (c *Charge) SetID(id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	c.id = id
	return nil
}

func (c *Charge) SetUserID(userID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("userID must not be nil")
	}
	c.userID = userID
	return nil
}

func (c *Charge) SetAmount(amount int) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	c.amount = amount
	return nil
}

func (c *Charge) SetBalance(balance int) error {
	if balance < 0 {
		return errors.New("balance must be non-negative")
	}
	c.balance = balance
	return nil
}

func (c *Charge) SetCreatedAt(createdAt time.Time) error {

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

	c.createdAt = createdAt

	return nil
}

func (c *Charge) SetDeletedAt(deletedAt time.Time) error {

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

	c.deletedAt = deletedAt

	return nil
}

func (c *Charge) ID() int {
	return c.id
}

func (c *Charge) UserID() uuid.UUID {
	return c.userID
}

func (c *Charge) Amount() int {
	return c.amount
}

func (c *Charge) Balance() int {
	return c.balance
}

func (c *Charge) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Charge) DeletedAt() time.Time {
	return c.deletedAt
}
