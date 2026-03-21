package purchase

import (
	"errors"
	"time"
)

type Purchase struct {
	id        int
	userID    string
	balance   int
	items     []PurchaseItem
	createdAt time.Time
	deletedAt time.Time
}

type PurchaseItem struct {
	itemID   int
	quantity int
}

func NewPurchaseItem(itemID int, quantity int) (PurchaseItem, error) {
	if itemID <= 0 {
		return PurchaseItem{}, errors.New("itemID must be positive")
	}
	if quantity <= 0 {
		return PurchaseItem{}, errors.New("quantity must be positive")
	}
	return PurchaseItem{itemID: itemID, quantity: quantity}, nil
}

func (i PurchaseItem) ItemID() int {
	return i.itemID
}

func (i PurchaseItem) Quantity() int {
	return i.quantity
}

func (p *Purchase) SetID(id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	p.id = id
	return nil
}

func (p *Purchase) SetUserID(userID string) error {
	if userID == "" {
		return errors.New("userID must not be empty")
	}
	p.userID = userID
	return nil
}

func (p *Purchase) SetBalance(balance int) error {
	if balance < 0 {
		return errors.New("balance must be non-negative")
	}
	p.balance = balance
	return nil
}

func (p *Purchase) SetItems(items []PurchaseItem) error {
	if len(items) == 0 {
		return errors.New("items must not be empty")
	}
	for _, item := range items {
		if item.itemID <= 0 {
			return errors.New("itemID must be positive")
		}
		if item.quantity <= 0 {
			return errors.New("quantity must be positive")
		}
	}
	p.items = items
	return nil
}

func (p *Purchase) SetCancelItems(items []PurchaseItem) error {
	if len(items) == 0 {
		return errors.New("items must not be empty")
	}
	for _, item := range items {
		if item.itemID <= 0 {
			return errors.New("itemID must be positive")
		}
		if item.quantity <= 0 {
			return errors.New("quantity must be positive")
		}
	}
	p.items = items
	return nil
}

func (p *Purchase) SetCreatedAt(createdAt time.Time) error {
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
	p.createdAt = createdAt
	return nil
}

func (p *Purchase) SetDeletedAt(deletedAt time.Time) error {
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
	p.deletedAt = deletedAt
	return nil
}

func (p *Purchase) ID() int {
	return p.id
}

func (p *Purchase) UserID() string {
	return p.userID
}

func (p *Purchase) Balance() int {
	return p.balance
}

func (p *Purchase) Items() []PurchaseItem {
	return p.items
}

func (p *Purchase) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Purchase) DeletedAt() time.Time {
	return p.deletedAt
}
