package charge

import (
	"github.com/google/uuid"
)

func NewCharge(userID uuid.UUID, amount int) (*Charge, error) {
	c := &Charge{}

	if err := c.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := c.SetAmount(amount); err != nil {
		return nil, err
	}

	return c, nil
}
