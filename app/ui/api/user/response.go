package user

import (
	"sysken-pay-api/app/domain/object/user"

	"github.com/google/uuid"
)

type PostUserResponse struct {
	Status    string    `json:"status"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt string    `json:"created_at"`
}

func toPostUserResponse(user *user.User) *PostUserResponse {
	return &PostUserResponse{
		Status:    "success",
		UserID:    user.ID(),
		UserName:  user.UserName(),
		CreatedAt: user.CreatedAt().Format("2006-01-02T15:04:05.000Z"),
	}
}
