package item

import "sysken-pay-api/app/domain/object/item"

// PostItemResponse Response: 商品登録のレスポンス
type PostItemResponse struct {
	Status    string `json:"status"`
	Id        int    `json:"id"`
	JanCode   string `json:"jan_code"`
	ItemName  string `json:"item_name"`
	Price     int    `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toPostItemResponse(item *item.Item) *PostItemResponse {
	return &PostItemResponse{
		Status:    "success",
		Id:        item.ID(),
		JanCode:   item.JanCode(),
		ItemName:  item.Name(),
		Price:     item.Price(),
		CreatedAt: item.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt: item.UpdatedAt().Format("2006-01-02T15:04:05.000Z"),
	}
}

// UpdateItemResponse Response: 商品更新のレスポンス
func toUpdateItemResponse(item *item.Item) *PostItemResponse {
	return &PostItemResponse{
		Status:    "success",
		Id:        item.ID(),
		JanCode:   item.JanCode(),
		ItemName:  item.Name(),
		Price:     item.Price(),
		CreatedAt: item.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt: item.UpdatedAt().Format("2006-01-02T15:04:05.000Z"),
	}
}

// GetItemResponse Response: 商品取得のレスポンス
type GetItemResponse struct {
	Status   string `json:"status"`
	ItemId   int    `json:"item_id"`
	JanCode  string `json:"jan_code"`
	ItemName string `json:"item_name"`
	Price    int    `json:"price"`
}

func toGetItemResponse(item *item.Item) *GetItemResponse {
	return &GetItemResponse{
		Status:   "success",
		ItemId:   item.ID(),
		JanCode:  item.JanCode(),
		ItemName: item.Name(),
		Price:    item.Price(),
	}
}

// GetAllItemsResponse Response: 商品取得のレスポンス
type GetAllItemResponse struct {
	ItemId   int    `json:"item_id"`
	JanCode  string `json:"jan_code"`
	ItemName string `json:"item_name"`
	Price    int    `json:"price"`
}

func toGetAllItemResponse(item *item.Item) *GetAllItemResponse {
	return &GetAllItemResponse{
		ItemId:   item.ID(),
		JanCode:  item.JanCode(),
		ItemName: item.Name(),
		Price:    item.Price(),
	}
}

// GetAllItemsResponse Response: 全商品取得のレスポンス
type GetAllItemsResponse struct {
	Status string               `json:"status"`
	Items  []GetAllItemResponse `json:"items"`
}

func toGetAllItemsResponse(items []*item.Item) *GetAllItemsResponse {
	itemResponses := make([]GetAllItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, *toGetAllItemResponse(item))
	}
	return &GetAllItemsResponse{
		Status: "success",
		Items:  itemResponses,
	}
}
