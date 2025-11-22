package item

import (
	"time"
)

func UpdateItem(janCode string, name string, price int) (*Item, error) {
	p := &Item{}

	if err := p.SetJanCode(janCode); err != nil {
		return nil, err
	}
	if err := p.SetName(name); err != nil {
		return nil, err
	}
	if err := p.SetPrice(price); err != nil {
		return nil, err
	}
	p.updatedAt = time.Now()

	return p, nil
}
