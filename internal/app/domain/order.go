package domain

import "time"

type RequestStoreOrderDTO struct {
	JSON string `json:"JSON,omitempty"`
}

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

	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`
}

type OrderService interface {
	GetByUID(int) (Order, error)
	StoreOrder(dto RequestStoreOrderDTO) (int, error)
}

type OrderRepository interface {
	Get(int) (Order, error)
	Store(RequestStoreOrderDTO) (int, error)
}
