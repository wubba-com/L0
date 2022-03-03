package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/client/pg"
	"log"
)

const (
	table = "orders"
)

func NewOrderRepository(client pg.Client) domain.OrderRepository {
	return &orderRepo{db: client}
}

type orderRepo struct {
	db pg.Client
}

var orderNotFound = errors.New("order not found")

func (r *orderRepo) CheckUnique(ctx context.Context, uid string) error {
	order := &domain.Order{}
	var query = fmt.Sprintf("SELECT order_uid FROM %s WHERE order_uid = $1", table)

	if err := r.db.QueryRow(ctx, query, &uid).Scan(&order.OrderUID); err != nil {
		return orderNotFound
	}

	return nil
}

// Get - get order with links
func (r *orderRepo) Get(ctx context.Context, uid string) (*domain.Order, error) {
	err := r.CheckUnique(context.Background(), uid)
	if err != nil {
		return nil, err
	}

	var query1 = fmt.Sprintf(
		"SELECT " +
			"orders.order_uid, " +
			"orders.track_number, " +
			"orders.entry, " +
			"deliveries.order_uid, " +
			"deliveries.name, " +
			"deliveries.phone, " +
			"deliveries.zip, " +
			"deliveries.city, " +
			"deliveries.address, " +
			"deliveries.region, " +
			"deliveries.email, " +
			"payments.transaction, " +
			"payments.request_id, " +
			"payments.currency, " +
			"payments.provider, " +
			"payments.amount, " +
			"payments.payment_dt, " +
			"payments.bank, " +
			"payments.delivery_cost, " +
			"payments.goods_total, " +
			"payments.custom_fee, " +
			"orders.locale, " +
			"orders.internal_signature, " +
			"orders.customer_id, " +
			"orders.delivery_service, " +
			"orders.shardkey, " +
			"orders.sm_id, " +
			"orders.date_created, " +
			"orders.oof_shard " +
			"FROM orders " +
			"JOIN deliveries ON deliveries.order_uid = orders.order_uid " +
			"JOIN payments ON payments.transaction = orders.order_uid " +
			"WHERE orders.order_uid = $1;")

	// init structs
	order := &domain.Order{}
	delivery := &domain.Delivery{}
	payment := &domain.Payment{}

	// get delivery and payment by order_uid
	row := r.db.QueryRow(ctx, query1, &uid)
	err = row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&delivery.OrderUID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
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

	// added full structs
	order.Delivery = delivery
	order.Payment = payment

	// get items by order_uid
	var query2 = fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1")
	rows, err := r.db.Query(ctx, query2, &order.OrderUID)
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

// All - get all orders
func (r *orderRepo) All(ctx context.Context) ([]*domain.Order, error) {
	var query = fmt.Sprintf(
		"SELECT "+
			"orders.order_uid,"+
			"orders.track_number, "+
			"orders.entry, "+
			"deliveries.order_uid, "+
			"deliveries.name, "+
			"deliveries.phone, "+
			"deliveries.zip, "+
			"deliveries.city, "+
			"deliveries.address, "+
			"deliveries.region, "+
			"deliveries.email, "+
			"payments.transaction, "+
			"payments.request_id, "+
			"payments.currency, "+
			"payments.provider, "+
			"payments.amount, "+
			"payments.payment_dt, "+
			"payments.bank, "+
			"payments.delivery_cost, "+
			"payments.goods_total, "+
			"payments.custom_fee, "+
			"orders.locale, "+
			"orders.internal_signature, "+
			"orders.customer_id, "+
			"orders.delivery_service, "+
			"orders.shardkey, "+
			"orders.sm_id, "+
			"orders.date_created, "+
			"orders.oof_shard "+
			"FROM %s "+
			"JOIN deliveries ON deliveries.order_uid = orders.order_uid "+
			"JOIN payments ON payments.transaction = orders.order_uid "+
			"ORDER BY orders.date_created LIMIT 10000", table)

	// get all order with links delivery, payment
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	orders := make([]*domain.Order, 0)
	for rows.Next() {

		// init structs
		order := &domain.Order{}
		delivery := &domain.Delivery{}
		payment := &domain.Payment{}

		if err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&delivery.OrderUID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDt,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
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

		// get link item by order
		var query2 = fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1")
		rows, err := r.db.Query(ctx, query2, &order.OrderUID)
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

		// added full structs
		order.Delivery = delivery
		order.Payment = payment
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *orderRepo) Store(ctx context.Context, order *domain.Order) (string, error) {
	var orderUID string
	query := fmt.Sprintf("INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid")
	if err := r.db.QueryRow(
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
			fmt.Println(fmt.Errorf("SQL order: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
			return "", err
		}
		log.Printf("[err] db: %s\n", err.Error())
		return "", err
	}

	return orderUID, nil
}
