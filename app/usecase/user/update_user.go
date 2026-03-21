package user

import (
	"context"
	"log"
	"sysken-pay-api/app/domain/object/user"
	"sysken-pay-api/app/domain/repository"
)

// TODO ドメイン層のインターフェースに接続をして処理を完成させる
// ユースケースとして、APIから受け取ったデータをドメイン層に渡す役割を果たす

type UpdateUserUseCase interface {
	UpdateUser(ctx context.Context, userId string, userName string) (*user.User, error)
}

type UpdateUserServiceImpl struct {
	userUpdateRepo repository.UserRepository
}

func NewUpdateUserUseCase(
	userUpdateRepo repository.UserRepository,
) *UpdateUserServiceImpl {
	return &UpdateUserServiceImpl{
		userUpdateRepo: userUpdateRepo,
	}
}

func (s *UpdateUserServiceImpl) UpdateUser(
	ctx context.Context, userId string, userName string) (*user.User, error) {

	// ドメインオブジェクトの生成
	u, err := user.UpdateUser(userName)
	if err != nil {
		log.Printf("Failed to create user object: %v", err)
		return nil, err
	}
	// IDをセット（UpdateUserファクトリでは設定されないため）
	if err := u.SetUserID(userId); err != nil {
		log.Printf("Failed to set user id: %v", err)
		return nil, err
	}

	updatedUser, err := s.userUpdateRepo.UpdateUser(ctx, u)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return updatedUser, nil
}
