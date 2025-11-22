package item

func GetItemByJanCode(janCode string) (*Item, error) {
	p := &Item{}

	if err := p.SetJanCode(janCode); err != nil {
		return nil, err
	}

	return p, nil
}
