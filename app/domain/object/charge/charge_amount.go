package charge

func NewCharge(userID string, amount int) (*Charge, error) {
	c := &Charge{}

	if err := c.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := c.SetAmount(amount); err != nil {
		return nil, err
	}

	return c, nil
}
