package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/client/pg"
	"log"
)

const(
	table = "orders"
)

func NewOrderRepository(client pg.Client) domain.OrderRepository {
	return &repository{p: client}
}

type repository struct {
	p pg.Client
}

func (r *repository) Get(ctx context.Context, id string) (*domain.Order, error) {
	//	"SELECT " +
	//		"order_uid, " +
	//		"track_number, " +
	//		"entry, " +
	//		"locale, " +
	//		"internal_signature, " +
	//		"customer_id, " +
	//		"delivery_service, " +
	//		"shardkey, " +
	//		"sm_id, " +
	//		"date_created, " +
	//		"oof_shard " +
	//		"FROM orders WHERE order_uid = $1")

	var query = fmt.Sprintf(
		`SELECT
			   orders.order_uid,
			   orders.track_number,
			   orders.entry,
			   deliveries.order_uid,
			   deliveries.name,
			   deliveries.phone,
			   deliveries.zip,
			   deliveries.city,
			   deliveries.address,
			   deliveries.region,
			   deliveries.email,
			   payments.transaction,
			   payments.request_id,
			   payments.currency,
			   payments.provider,
			   payments.amount,
			   payments.payment_dt,
			   payments.bank,
			   payments.delivery_cost,
			   payments.goods_total,
			   payments.custom_fee,
			   orders.locale,
			   orders.internal_signature,
			   orders.customer_id,
			   orders.delivery_service,
			   orders.shardkey,
			   orders.sm_id,
			   orders.date_created,
			   orders.oof_shard
		FROM orders WHERE order_uid = $1
		JOIN deliveries ON deliveries.order_uid = orders.order_uid 
		JOIN payments ON payments.order_uid = orders.order_uid;`) //query := fmt.Sprintf(

	order := &domain.Order{}

	err := r.p.QueryRow(ctx, query, &id).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		)
	if err != nil {
		log.Printf("[err] repository:%s\n", err.Error())
		return nil, err
	}

	var query2 = fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1")
	rows, err := r.p.Query(ctx, query2, &order.OrderUID)
	if err != nil {
		return nil, err
	}

	items := make([]*domain.Item, 0)
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
	order.Items = items
	return order, nil
}

func (r *repository) All(ctx context.Context) ([]*domain.Order, error) {
	query := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,sm_id, date_created, oof_shard FROM %s ORDER BY DESK date_created LIMIT 10000", table)
	rows, err := r.p.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		if err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			); err != nil {

			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) Store(ctx context.Context, order *domain.Order) (string, error) {
	var orderUID string
	query := fmt.Sprintf("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid")
	if err := r.p.QueryRow(
		ctx,
		query,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	).Scan(&orderUID); err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Errorf("SQL: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
			return "", err
		}
		log.Printf("[err] db: %s\n", err.Error())
		return "", err
	}

	return orderUID, nil
}
