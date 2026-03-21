package purchase

import "errors"

// CanPurchase はユーザーの残高が購入金額以上かを検証するドメインサービス
func CanPurchase(balance, totalAmount int) error {
	if balance < totalAmount {
		return errors.New("insufficient balance")
	}
	return nil
}
