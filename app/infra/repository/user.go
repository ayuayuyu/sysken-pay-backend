package repository

import (
	"context"
	"database/sql"
	"log"
	"sysken-pay-api/app/domain/object/user"
	"sysken-pay-api/app/domain/repository"
	"time"
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
	ctx context.Context, u *user.User) (*user.User, error) {

	executor := getExecutor(ctx, r.db)

	query := `
    INSERT INTO ` + "`user`" + ` (id, name, deleted_at)
    VALUES (?, ?,  NULL)
	`
	_, err := executor.ExecContext(ctx, query,
		u.ID(),
		u.UserName(),
	)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	row := executor.QueryRowContext(ctx, `
    SELECT created_at, updated_at FROM `+"`user`"+` WHERE id = ?
	`, u.ID())

	var createdAt, updatedAt time.Time
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		return nil, err
	}

	u.SetCreatedAt(createdAt)
	u.SetUpdatedAt(updatedAt)

	return u, nil
}

func (r *UserRepositoryImpl) UpdateUser(
	ctx context.Context, u *user.User) (*user.User, error) {

	executor := getExecutor(ctx, r.db)

	query := `
	UPDATE ` + "`user`" + ` SET name = ? WHERE id = ? AND deleted_at IS NULL
	`
	_, err := executor.ExecContext(ctx, query,
		u.UserName(),
		u.ID(),
	)

	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}
	row := executor.QueryRowContext(ctx, `
	SELECT created_at, updated_at FROM `+"`user`"+` WHERE id = ?
	`, u.ID())

	var createdAt, updatedAt time.Time
	if err := row.Scan(&createdAt, &updatedAt); err != nil {
		return nil, err
	}

	u.SetCreatedAt(createdAt)
	u.SetUpdatedAt(updatedAt)

	return u, nil
}
