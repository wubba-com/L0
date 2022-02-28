package domain

import (
	"context"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid,omitempty"`
	TrackNumber       string    `json:"track_number,omitempty"`
	Entry             string    `json:"entry,omitempty"`
	Locale            string    `json:"locale,omitempty"`
	InternalSignature string    `json:"internal_signature,omitempty"`
	CustomerID        string    `json:"customer_id,omitempty"`
	DeliveryService   string    `json:"delivery_service,omitempty"`
	ShardKey          string    `json:"shardkey,omitempty"`
	SmID              uint64    `json:"sm_id,omitempty"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard,omitempty"`

	Delivery *Delivery `json:"delivery"`
	Payment  *Payment  `json:"payment"`
	Items    []*Item   `json:"items"`
}

type OrderService interface {
	GetByUID(context.Context, string) (*Order, error)
	StoreOrder(context.Context, *Order) (string, error)
	LoadOrderCache(ctx context.Context) error
}

type OrderRepository interface {
	Get(context.Context, string) (*Order, error)
	Store(context.Context, *Order) (string, error)
	All(context.Context) ([]*Order, error)
}
