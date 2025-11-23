package query

import (
	"context"
	"database/sql"
	"log"
	"sysken-pay-api/app/domain/object/user"
	"sysken-pay-api/app/domain/repository"
	"time"

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
    INSERT INTO ` + "`user`" + ` (id, name, deleted_at)
    VALUES (?, ?,  NULL)
	`
	_, err = r.db.ExecContext(ctx, query,
		u.ID(),
		u.UserName(),
	)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, `
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
	ctx context.Context, userID uuid.UUID, userName string) (*user.User, error) {

	u, err := user.UpdateUser(userName)

	if err != nil {
		log.Printf("Failed to create user domain object: %v", err)
		return nil, err
	}

	u.SetUserID(userID)

	query := `
	UPDATE ` + "`user`" + ` SET name = ? WHERE id = ? AND deleted_at IS NULL
	`
	_, err = r.db.ExecContext(ctx, query,
		u.UserName(),
		u.ID(),
	)

	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}
	row := r.db.QueryRowContext(ctx, `
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
