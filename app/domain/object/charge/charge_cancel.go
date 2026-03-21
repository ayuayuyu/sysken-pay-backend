package charge

func DeleteCharge(userId string, amount int) (*Charge, error) {
	c := &Charge{}

	if err := c.SetUserID(userId); err != nil {
		return nil, err
	}
	if err := c.SetAmount(-amount); err != nil {
		return nil, err
	}

	return c, nil
}
