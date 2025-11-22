package item

import (
	"time"
)

func NewItem(janCode string, name string, price int) (*Item, error) {
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
	p.createdAt = time.Now()
	p.updatedAt = time.Now()

	return p, nil
}
