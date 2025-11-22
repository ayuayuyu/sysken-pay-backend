package item

import (
	"errors"
	"time"
	"unicode/utf8"
)

type Item struct {
	id        int
	janCode   string
	name      string
	price     int
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func (p *Item) SetID(id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}
	p.id = id
	return nil
}

func (p *Item) SetJanCode(janCode string) error {
	if janCode == "" {
		return errors.New("janCode must not be empty")
	}
	// JANコードは13桁の数字か短縮型の8桁であること
	if len(janCode) != 13 && len(janCode) != 8 {
		return errors.New("janCode must be 13 or 8 characters long")
	}
	// JANコードは数字のみで構成されていること
	if !isNumeric(janCode) {
		return errors.New("janCode must be numeric")
	}
	// JANコードのチェックデジットが正しいこと
	if !isValidJanCode(janCode) {
		return errors.New("janCode must be valid")
	}
	p.janCode = janCode
	return nil
}

func (p *Item) SetName(name string) error {
	if !(utf8.RuneCountInString(name) >= 1) {
		return errors.New("name must not be empty")
	}
	p.name = name
	return nil
}

func (p *Item) SetPrice(price int) error {
	if price < 0 {
		return errors.New("price must be non-negative")
	}
	p.price = price
	return nil
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
func isValidJanCode(janCode string) bool {
	length := len(janCode)
	checkDigit := int(janCode[length-1] - '0')
	sum := 0
	multiplier := 3
	for i := length - 2; i >= 0; i-- {
		digit := int(janCode[i] - '0')
		sum += digit * multiplier
		if multiplier == 3 {
			multiplier = 1
		} else {
			multiplier = 3
		}
	}
	calculatedCheckDigit := (10 - (sum % 10)) % 10
	return checkDigit == calculatedCheckDigit
}

func (p *Item) ID() int {
	return p.id
}

func (p *Item) JanCode() string {
	return p.janCode
}
func (p *Item) Name() string {
	return p.name
}
func (p *Item) Price() int {
	return p.price
}
func (p *Item) CreatedAt() time.Time {
	return p.createdAt
}
func (p *Item) UpdatedAt() time.Time {
	return p.updatedAt
}
func (p *Item) DeletedAt() time.Time {
	return p.deletedAt
}

func NewItemFromDB(
	id int,
	janCode string,
	name string,
	price int,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt time.Time,
) (*Item, error) {
	return &Item{
		id:        id,
		janCode:   janCode,
		name:      name,
		price:     price,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}, nil
}
