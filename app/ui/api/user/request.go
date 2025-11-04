package user

// PostUser Request: ユーザー新規登録のリクエスト
type PostUserRequest struct {
	UserName string `json:"user_name"`
}
