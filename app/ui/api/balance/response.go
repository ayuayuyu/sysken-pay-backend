package balance

import (
	domainbalance "sysken-pay-api/app/domain/object/balance"
	"sysken-pay-api/app/usecase/balance"
)

type GetBalanceResponse struct {
	Status  string `json:"status"`
	UserID  string `json:"user_id"`
	Balance int    `json:"balance"`
}

func toGetBalanceResponse(b *domainbalance.Balance) *GetBalanceResponse {
	return &GetBalanceResponse{
		Status:  "success",
		UserID:  b.UserID(),
		Balance: b.Balance(),
	}
}

type HistoryItem struct {
	ItemName   string `json:"item_name"`
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	PurchaseAt string `json:"purchase_at"`
}

type PageInfo struct {
	PrevPage  *int `json:"prev_page,omitempty"`
	NextPage  *int `json:"next_page,omitempty"`
	TotalPage int  `json:"total_page"`
}

type GetPurchaseHistoriesResponse struct {
	Status    string        `json:"status"`
	Histories []HistoryItem `json:"histories"`
	Total     int           `json:"total"`
	Page      PageInfo      `json:"page"`
}

func toGetPurchaseHistoriesResponse(
	histories []*domainbalance.PurchaseHistory,
	totalAmount int,
	page *balance.HistoryPage,
) *GetPurchaseHistoriesResponse {
	items := make([]HistoryItem, len(histories))
	for i, h := range histories {
		items[i] = HistoryItem{
			ItemName:   h.ItemName(),
			Quantity:   h.Quantity(),
			Price:      h.Price(),
			PurchaseAt: h.PurchaseAt().Format("2006-01-02T15:04:05.000Z"),
		}
	}
	return &GetPurchaseHistoriesResponse{
		Status:    "success",
		Histories: items,
		Total:     totalAmount,
		Page: PageInfo{
			PrevPage:  page.PrevPage(),
			NextPage:  page.NextPage(),
			TotalPage: page.TotalPage(),
		},
	}
}
