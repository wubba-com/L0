package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/pg"
)

func NewOrderRepository(client pg.Client) domain.OrderRepository {
	return &repository{}
}

type repository struct {
	p pg.Client
}

func (r repository) Get(ctx context.Context, id int) (*domain.Order, error) {
	order := &domain.Order{}

	query := `SELECT * FROM orders WHERE order_uid = $1`

	err := r.p.QueryRow(ctx, query, id).Scan(order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r repository) Store(ctx context.Context, order *domain.Order) (string, error) {
	var orderUID string
	query := `INSERT INTO orders 
		(
			order_uid, 
			track_number, 
			entry, 
			locale, 
			internal_signature, 
			customer_id, 
			delivery_service, 
			shardkey, 
			sm_id, 
			date_created, 
			oof_shard
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING order_uid`
	if err := r.p.QueryRow(
		ctx,
		query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	).Scan(orderUID); err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok {
			println(fmt.Errorf("SQL: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
			return "", nil
		}
		return "", err
	}

	return orderUID, nil
}
