package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse は API のエラーレスポンスの共通フォーマットです。
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RespondError は指定されたステータスコードとメッセージで
// `{ "status": "error", "message": "..." }` 形式の JSON を返します。
func RespondError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// エンコードが失敗した場合はログ等へ出力する実装に置き換えてください。
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: msg,
	})
}
