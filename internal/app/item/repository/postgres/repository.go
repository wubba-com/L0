package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/client/pg"
)

const(
	table = "items"
)

func NewItemRepository(client pg.Client) domain.ItemRepository {
	return &itemRepo{p: client}
}

type itemRepo struct {
	p pg.Client
}

func (i *itemRepo) GetByOrderUID(ctx context.Context, orderUID string) ([]*domain.Item, error) {
	query := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid = $1", table)
	items := make([]*domain.Item, 0)

	rows, err := i.p.Query(ctx, query, &orderUID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &domain.Item{}
		if err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
			); err != nil {
				return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (i *itemRepo) Store(ctx context.Context, item *domain.Item) (uint64, error) {
	var ChrtID uint64
	var query = fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING chrt_id", table)
	if err := i.p.QueryRow(
		ctx,
		query,
		&item.ChrtID,
		&item.TrackNumber,
		&item.Price,
		&item.Rid,
		&item.Name,
		&item.Sale,
		&item.Size,
		&item.TotalPrice,
		&item.NmID,
		&item.Brand,
		&item.Status,
		&item.OrderUID,
		).Scan(&ChrtID); err != nil {
			if pgError, ok := err.(*pgconn.PgError); ok {
				fmt.Println(fmt.Errorf("SQL: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
				return 0, err
			}

		return 0, err
	}

	return  ChrtID, nil
}