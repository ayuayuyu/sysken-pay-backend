package charge

import (
	"errors"
	"time"
)

//TODO モデル（データベースに入れる型を宣言する）
//データベースの制約通りになるようにエラーハンドリングをガチる

type Charge struct {
	id        int
	userID    string
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

func (c *Charge) SetUserID(userID string) error {
	if userID == "" {
		return errors.New("userID must not be empty")
	}
	c.userID = userID
	return nil
}

// allowedAmounts はチャージ可能な固定金額の一覧
var allowedAmounts = map[int]struct{}{
	500:  {},
	1000: {},
	2000: {},
}

func (c *Charge) SetAmount(amount int) error {
	if _, ok := allowedAmounts[amount]; !ok {
		return errors.New("amount must be one of 500, 1000, 2000")
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

func (c *Charge) UserID() string {
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
