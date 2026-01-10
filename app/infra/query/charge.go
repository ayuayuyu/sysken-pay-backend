package query

import (
	"context"
	"database/sql"
	"sysken-pay-api/app/domain/object/charge"
	"sysken-pay-api/app/domain/repository"
	"time"

	"github.com/google/uuid"
)

var _ repository.ChargeRepository = (*ChargeRepositoryImpl)(nil)

type ChargeRepositoryImpl struct {
	db *sql.DB
}

func NewChargeRepository(db *sql.DB) *ChargeRepositoryImpl {
	return &ChargeRepositoryImpl{
		db: db,
	}
}

// TODO ChargeAmountメソッドの実装
func (r *ChargeRepositoryImpl) ChargeAmount(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {

	c, err := charge.NewCharge(userID, amount)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO charges (user_id, amount) VALUES (?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		c.UserID(),
		c.Amount(),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	c.SetID(int(id))

	row := r.db.QueryRowContext(
		ctx,
		`SELECT created_at, updated_at FROM charges WHERE id = ?`,
		c.ID(),
	)

	var createdAt time.Time
	if err := row.Scan(&createdAt); err != nil {
		return nil, err
	}

	c.SetCreatedAt(createdAt)

	return c, nil
}

// TODO ChargeCancelメソッドの実装
func (r *ChargeRepositoryImpl) ChargeCancel(ctx context.Context, userID uuid.UUID, amount int) (*charge.Charge, error) {
	c, err := charge.DeleteCharge(userID, amount)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO charges (user_id, amount) VALUES (?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		c.UserID(),
		c.Amount(),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	c.SetID(int(id))

	row := r.db.QueryRowContext(
		ctx,
		`SELECT created_at, updated_at FROM charges WHERE id = ?`,
		c.ID(),
	)

	var deletedAt time.Time
	if err := row.Scan(&deletedAt); err != nil {
		return nil, err
	}

	c.SetDeletedAt(deletedAt)

	return c, nil
}
