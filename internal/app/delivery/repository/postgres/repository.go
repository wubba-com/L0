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
	table = "deliveries"
)

func NewDeliveryRepository(client pg.Client) domain.DeliveryRepository {
	return &deliveryRepo{client}
}

type deliveryRepo struct {
	db pg.Client
}

func (d *deliveryRepo) GetByOrderUID(ctx context.Context, orderUID string) (*domain.Delivery, error) {
	query := fmt.Sprintf("SELECT delivery_uid, order_uid, name, phone, zip, city, address, region, email FROM %s WHERE order_uid = $1", table)
	delivery := &domain.Delivery{}

	if err := d.db.QueryRow(
		ctx, query, orderUID).Scan(
			&delivery.DeliveryUID,
			&delivery.OrderUID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		); err != nil {
		return nil, err
	}

	return delivery, nil
}

func (d *deliveryRepo) Store(ctx context.Context, delivery *domain.Delivery) (uint64, error) {
	var deliveryUID uint64
	query := fmt.Sprintf("INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $6, $7, $8) RETURNING delivery_id")
	if err := d.db.QueryRow(
		ctx,
		query,
		&delivery.OrderUID,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	).Scan(&deliveryUID); err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Errorf("SQL: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
			return 0, err
		}
		log.Printf("[err] db: %s\n", err.Error())
		return 0, err
	}
	return deliveryUID, nil
}
