package purchase

func DeletePurchase(userId string, items []PurchaseItem) (*Purchase, error) {
	p := &Purchase{}

	if err := p.SetUserID(userId); err != nil {
		return nil, err
	}
	if err := p.SetItems(items); err != nil {
		return nil, err
	}

	return p, nil
}
