package query

import (
	"context"
	"database/sql"
	"sysken-pay-api/app/domain/object/item"
	"sysken-pay-api/app/domain/repository"
	"time"
)

//TODO itemデータベースに値を入れる
//domainのrepositoryの中にあるアイテムのインタフェースを実装する

var _ repository.ItemRepository = (*ItemRepositoryImpl)(nil)

type ItemRepositoryImpl struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepositoryImpl {
	return &ItemRepositoryImpl{
		db: db,
	}
}

// TODO InsertItemメソッドの実装
func (r *ItemRepositoryImpl) InsertItem(ctx context.Context, janCode string, name string, price int) (*item.Item, error) {
	i, err := item.NewItem(janCode, name, price)
	if err != nil {
		return nil, err
	}

	query := `
	INSERT INTO item (jan_code, name, price, created_at,updated_at, deleted_at)
	VALUES (?, ?, ?, ?, ?, NULL)
	`
	result, err := r.db.ExecContext(ctx, query,
		i.JanCode(),
		i.Name(),
		i.Price(),
		i.CreatedAt(),
		i.UpdatedAt(),
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	i.SetID(int(id))

	return i, nil
}

// TODO UpdateItemメソッドの実装
func (r *ItemRepositoryImpl) UpdateItem(ctx context.Context, janCode string, name string, price int) (*item.Item, error) {
	i, err := item.UpdateItem(janCode, name, price)
	if err != nil {
		return nil, err
	}

	query := `
	UPDATE item
	SET name = ?, price = ?, updated_at = ?
	WHERE jan_code = ? AND deleted_at IS NULL
	`
	_, err = r.db.ExecContext(ctx, query,
		i.Name(),
		i.Price(),
		i.UpdatedAt(),
		i.JanCode(),
	)

	if err != nil {
		return nil, err
	}

	var id int
	err = r.db.QueryRowContext(
		ctx,
		`SELECT id FROM item WHERE jan_code = ? AND deleted_at IS NULL`,
		i.JanCode(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	i.SetID(id)

	return i, nil
}

// TODO GetItemByJanCodeメソッドの実装
func (r *ItemRepositoryImpl) GetItemByJanCode(ctx context.Context, janCode string) (*item.Item, error) {

	query := `
	SELECT id, jan_code, name, price, created_at, updated_at, deleted_at
	FROM item
	WHERE jan_code = ? AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, janCode)

	var (
		id        int
		janCodeDB string
		nameDB    string
		priceDB   int
		createdAt time.Time
		updatedAt time.Time
		deletedAt sql.NullTime
	)

	if err := row.Scan(
		&id,
		&janCodeDB,
		&nameDB,
		&priceDB,
		&createdAt,
		&updatedAt,
		&deletedAt,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// deletedAt が NULL ならゼロ値 time.Time にする
	var deleted time.Time
	if deletedAt.Valid {
		deleted = deletedAt.Time
	}

	i, err := item.NewItemFromDB(
		id,
		janCodeDB,
		nameDB,
		priceDB,
		createdAt,
		updatedAt,
		deleted,
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}

// TODO GetAllItemsメソッドの実装
func (r *ItemRepositoryImpl) GetAllItems(ctx context.Context) ([]*item.Item, error) {
	query := `
	SELECT id, jan_code, name, price, created_at, updated_at, deleted_at
	FROM item
	WHERE deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*item.Item

	for rows.Next() {
		var (
			id        int
			janCodeDB string
			nameDB    string
			priceDB   int
			createdAt time.Time
			updatedAt time.Time
			deletedAt sql.NullTime
		)

		if err := rows.Scan(
			&id,
			&janCodeDB,
			&nameDB,
			&priceDB,
			&createdAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		// deletedAt が NULL ならゼロ値 time.Time にする
		var deleted time.Time
		if deletedAt.Valid {
			deleted = deletedAt.Time
		}

		i, err := item.NewItemFromDB(
			id,
			janCodeDB,
			nameDB,
			priceDB,
			createdAt,
			updatedAt,
			deleted,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
