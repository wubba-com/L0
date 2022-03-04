package domain

import (
	"context"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid,omitempty" validate:"required"`
	TrackNumber       string    `json:"track_number,omitempty" validate:"required,max=128"`
	Entry             string    `json:"entry,omitempty" validate:"required,max=128"`
	Locale            string    `json:"locale,omitempty" validate:"required,max=128"`
	InternalSignature string    `json:"internal_signature,omitempty" validate:"max=128"`
	CustomerID        string    `json:"customer_id,omitempty" validate:"required,max=128"`
	DeliveryService   string    `json:"delivery_service,omitempty" validate:"required,max=128"`
	ShardKey          string    `json:"shardkey,omitempty" validate:"required,max=128"`
	SmID              uint64    `json:"sm_id,omitempty" validate:"required,gt=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard,omitempty" validate:"required,max=128"`

	Delivery *Delivery `json:"delivery"`
	Payment  *Payment  `json:"payment"`
	Items    []*Item   `json:"items" validate:"required,dive"`
}

type OrderService interface {
	GetByUID(context.Context, string) (*Order, error)
	StoreOrder(context.Context, *Order) (string, error)
	All(ctx context.Context) ([]*Order, error)
	LoadOrderCache(ctx context.Context) error
}

type OrderRepository interface {
	Get(context.Context, string) (*Order, error)
	Store(context.Context, *Order) (string, error)
	All(context.Context) ([]*Order, error)
	CheckUnique(context.Context, string) error
}
