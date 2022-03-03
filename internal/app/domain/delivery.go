package domain

import "context"

type Delivery struct {
	DeliveryUID uint64
	Name    string `json:"name,omitempty" validate:"required,max=128"`
	Phone   string `json:"phone,omitempty" validate:"required,max=32"`
	Zip     string `json:"zip,omitempty" validate:"required,max=128"`
	City    string `json:"city,omitempty" validate:"required,max=128"`
	Address string `json:"address,omitempty" validate:"required,max=128"`
	Region  string `json:"region,omitempty" validate:"required,max=128"`
	Email   string `json:"email,omitempty" validate:"required,max=128,email"`

	OrderUID string `validate:"required"`
}

type DeliveryRepository interface {
	GetByOrderUID(context.Context, string) (*Delivery, error)
	Store(context.Context, *Delivery) (uint64, error)
}

type DeliveryService interface {
	GetDeliveryByOrderUID(context.Context, string) (*Delivery, error)
	StoreDelivery(context.Context, *Delivery) (uint64, error)
}