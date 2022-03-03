package order

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/cache"
	"time"
)

func NewOrderService(repository domain.OrderRepository, c cache.Cache, ttl time.Duration) domain.OrderService {
	return &serviceOrder{r: repository, c: c, ttlCache: ttl}
}

type serviceOrder struct {
	r        domain.OrderRepository
	c        cache.Cache
	ttlCache time.Duration
}

func (s *serviceOrder) LoadOrderCache(ctx context.Context) ([]*domain.Order, error) {
	orders, err := s.r.All(ctx)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		s.c.Set(order.OrderUID, order, s.ttlCache)
	}
	return orders, nil
}

func (s *serviceOrder) GetByUID(ctx context.Context, uid string) (*domain.Order, error) {
	if order, found := s.c.Get(uid); found {
		return order.(*domain.Order), nil
	}
	order, err := s.r.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	s.c.Set(order.OrderUID, order, s.ttlCache)

	return order, nil
}

func (s *serviceOrder) StoreOrder(ctx context.Context, order *domain.Order) (string, error) {
	uid, err := s.r.Store(ctx, order)
	if err != nil {
		return "", err
	}
	s.c.Set(uid, order, s.ttlCache)

	return uid, nil
}
