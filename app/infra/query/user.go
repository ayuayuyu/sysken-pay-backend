package query

import (
	"context"
	"database/sql"
	"log"
	"sysken-pay-api/app/domain/object/user"
	"sysken-pay-api/app/domain/repository"

	"github.com/google/uuid"
)

//TODO userデータベースに値を入れる
// domainのrepositoryの中にあるユーザーのインターフェースの実装をする

var _ repository.UserRepository = (*RegisterUserServiceImpl)(nil)

type RegisterUserServiceImpl struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *RegisterUserServiceImpl {
	return &RegisterUserServiceImpl{db: db}
}

type InsertUserDTO struct {
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt string    `json:"created_at"`
	DeletedAt string    `json:"deleted_at"`
}

func (r *RegisterUserServiceImpl) InsertUser(
	ctx context.Context, userName string) (*user.User, error) {
	u, err := user.NewUser(uuid.New(), userName)

	if err != nil {
		log.Printf("Failed to create user domain object: %v", err)
		return nil, err
	}

	query := `
    INSERT INTO ` + "`user`" + ` (name, created_at, deleted_at)
    VALUES (?, ?, NULL)
	`
	_, err = r.db.ExecContext(ctx, query,
		u.UserName(),
		u.CreatedAt(),
	)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return u, nil
}
