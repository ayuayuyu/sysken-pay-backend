package balance

import (
	"errors"
	"time"
)

// PurchaseHistory は購入履歴の1件を表すドメインオブジェクト
type PurchaseHistory struct {
	itemName   string
	quantity   int
	price      int
	purchaseAt time.Time
}

func NewPurchaseHistory(itemName string, quantity, price int, purchaseAt time.Time) (*PurchaseHistory, error) {
	if itemName == "" {
		return nil, errors.New("itemName must not be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	if price < 0 {
		return nil, errors.New("price must be non-negative")
	}
	return &PurchaseHistory{
		itemName:   itemName,
		quantity:   quantity,
		price:      price,
		purchaseAt: purchaseAt,
	}, nil
}

func (h *PurchaseHistory) ItemName() string      { return h.itemName }
func (h *PurchaseHistory) Quantity() int         { return h.quantity }
func (h *PurchaseHistory) Price() int            { return h.price }
func (h *PurchaseHistory) PurchaseAt() time.Time { return h.purchaseAt }
