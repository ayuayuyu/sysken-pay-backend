package user

import (
	"context"
	"log"
	"sysken-pay-api/app/domain/object/user"
	"sysken-pay-api/app/domain/repository"
)

// TODO ドメイン層のインターフェースに接続をして処理を完成させる
// ユースケースとして、APIから受け取ったデータをドメイン層に渡す役割を果たす

type RegisterUserUseCase interface {
	RegisterUser(ctx context.Context, userName string) (*user.User, error)
}

type RegisterUserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewRegisterUserUseCase(
	userRepo repository.UserRepository,
) *RegisterUserServiceImpl {
	return &RegisterUserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *RegisterUserServiceImpl) RegisterUser(
	ctx context.Context, userName string) (*user.User, error) {

	createdUser, err := s.userRepo.InsertUser(ctx, userName)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return createdUser, nil
}
