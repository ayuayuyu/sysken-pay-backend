package user

//TODO モデル（データベースに入れる型を宣言する）
//データベースの制約通りになるようにエラーハンドリングをガチる
//ユーザーID、名前、作成日時、更新日時など

func NewUser(userID string, userName string) (*User, error) {
	user := &User{}

	if err := user.SetUserID(userID); err != nil {
		return nil, err
	}
	if err := user.SetUserName(userName); err != nil {
		return nil, err
	}
	return user, nil
}
