package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/client/pg"
)

const (
	table = "payments"
)

func NewPaymentRepository(client pg.Client) domain.PaymentRepository {
	return &paymentRepo{db: client}
}

type paymentRepo struct {
	db pg.Client
}

func (p *paymentRepo) GetByOrderUID(ctx context.Context, orderUID string) (*domain.Payment, error) {
	query := fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM %s WHERE transaction = $1", table)
	payment := &domain.Payment{}

	if err := p.db.QueryRow(
		ctx, query, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	); err != nil {
		return nil, err
	}

	return payment, nil
}

func (p *paymentRepo) Store(ctx context.Context, payment *domain.Payment) (string, error) {
	var Transaction string
	var query = fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING transaction", table)
	if err := p.db.QueryRow(
		ctx,
		query,
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
	).Scan(&Transaction); err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok {
			fmt.Println(fmt.Errorf("SQL payment: Error: %s, Detail:%s, Where: %s, Code:%s", pgError.Message, pgError.Detail, pgError.Where, pgError.Code))
			return "", err
		}

		return "", err
	}

	return Transaction, nil
}
