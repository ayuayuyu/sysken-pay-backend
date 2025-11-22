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

var _ repository.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) InsertUser(
	ctx context.Context, userName string) (*user.User, error) {
	u, err := user.NewUser(uuid.New(), userName)

	if err != nil {
		log.Printf("Failed to create user domain object: %v", err)
		return nil, err
	}

	query := `
    INSERT INTO ` + "`user`" + ` (id, name, created_at, deleted_at)
    VALUES (?, ?, ?, NULL)
	`
	_, err = r.db.ExecContext(ctx, query,
		u.ID(),
		u.UserName(),
		u.CreatedAt(),
	)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return u, nil
}
