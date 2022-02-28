package order

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
)

func NewOrderService(repository domain.OrderRepository) domain.OrderService {
	return &service{r: repository}
}

type service struct {
	r domain.OrderRepository
}

func (s service) GetByUID(ctx context.Context, id int) (*domain.Order, error) {
	order, err := s.r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s service) StoreOrder(ctx context.Context, order *domain.Order) (string, error) {
	id, err := s.r.Store(ctx, order)
	if err != nil {
		return "", err
	}

	return id, nil
}
