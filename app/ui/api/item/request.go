package item

// PostItem Request: 商品新規登録と更新のリクエスト
type PostItemRequest struct {
	JanCode  string `json:"jan_code"`
	ItemName string `json:"item_name"`
	Price    int    `json:"price"`
}
