package user

// PostUser Request: ユーザー新規登録のリクエスト
type PostUserRequest struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

// PatchUser Request: ユーザー更新のリクエスト
type PatchUserRequest struct {
	UserName string `json:"user_name"`
}
