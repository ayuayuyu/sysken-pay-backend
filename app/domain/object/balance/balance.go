package balance

import "errors"

// Balance はユーザーの現在の残高を表すドメインオブジェクト
type Balance struct {
	userID  string
	balance int
}

func NewBalance(userID string, balance int) (*Balance, error) {
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	if balance < 0 {
		return nil, errors.New("balance must be non-negative")
	}
	return &Balance{userID: userID, balance: balance}, nil
}

func (b *Balance) UserID() string { return b.userID }
func (b *Balance) Balance() int   { return b.balance }
