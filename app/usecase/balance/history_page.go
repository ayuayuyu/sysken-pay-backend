package balance

import "errors"

// HistoryPage はページネーション情報を保持する
type HistoryPage struct {
	prevPage  *int
	nextPage  *int
	totalPage int
}

func NewHistoryPage(prevPage, nextPage *int, totalPage int) (*HistoryPage, error) {
	if totalPage < 0 {
		return nil, errors.New("totalPage must be non-negative")
	}
	return &HistoryPage{prevPage: prevPage, nextPage: nextPage, totalPage: totalPage}, nil
}

func (p *HistoryPage) PrevPage() *int { return p.prevPage }
func (p *HistoryPage) NextPage() *int { return p.nextPage }
func (p *HistoryPage) TotalPage() int { return p.totalPage }
