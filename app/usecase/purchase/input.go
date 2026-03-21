package purchase

// PurchaseItemInput は購入・キャンセル時に UI 層から受け取る入力値
type PurchaseItemInput struct {
	ItemID   int
	Quantity int
}
