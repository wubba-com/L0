package domain

import "context"

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id" validate:"omitempty"`
	Currency     string `json:"currency" validate:"required,max=128"`
	Provider     string `json:"provider" validate:"required,max=128"`
	Amount       uint64 `json:"amount" validate:"required,gte=0"`
	PaymentDt    uint64 `json:"payment_dt" validate:"required,gte=0"`
	Bank         string `json:"bank" validate:"required,max=64"`
	DeliveryCost uint64 `json:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   uint64 `json:"goods_total" validate:"required,gte=0"`
	CustomFee    uint64 `json:"custom_fee" validate:"gte=0"`
}

type PaymentRepository interface {
	GetByOrderUID(context.Context, string) (*Payment, error)
	Store(context.Context, *Payment) (string, error)
}

type PaymentService interface {
	GetPaymentByOrderUID(context.Context, string) (*Payment, error)
	StorePayment(context.Context, *Payment) (string, error)
}
