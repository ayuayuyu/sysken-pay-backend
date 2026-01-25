package purchase

import (
	"github.com/google/uuid"
)

func NewPurchase(userID uuid.UUID, items []PurchaseItem) (*Purchase, error) {
	p := &Purchase{}

	if err := p.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := p.SetItems(items); err != nil {
		return nil, err
	}

	return p, nil
}
