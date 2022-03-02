package file

import (
	"context"
	"fmt"
	"github.com/wubba-com/L0/internal/app/domain"
)

func NewOrderRepository() domain.OrderRepository {
	return &repository{}
}

type repository struct {

}

func (r repository) Get(ctx context.Context, s string) (*domain.Order, error) {
	panic("implement me")
}

func (r repository) Store(ctx context.Context, order *domain.Order) (string, error) {
	fmt.Println(order.OrderUID, order.Delivery)
	return order.OrderUID, nil
}

func (r repository) All(ctx context.Context) ([]*domain.Order, error) {
	panic("implement me")
}

