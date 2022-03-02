package order

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/cache"
	"log"
	"time"
)

func NewOrderService(repository domain.OrderRepository, c cache.Cache, ttl time.Duration) domain.OrderService {
	return &service{r: repository, c: c, ttlCache: ttl}
}

type service struct {
	r        domain.OrderRepository
	c        cache.Cache
	ttlCache time.Duration
}

func (s *service) LoadOrderCache(ctx context.Context) error {
	orders, err := s.r.All(ctx)
	if err != nil {
		log.Printf("[err] service:%s\n", err.Error())
		return err
	}
	for _, order := range orders {
		s.c.Set(order.OrderUID, order, s.ttlCache)
	}
	return nil
}

func (s *service) GetByUID(ctx context.Context, uid string) (*domain.Order, error) {
	if order, found := s.c.Get(uid); found {
		return order.(*domain.Order), nil
	}
	order, err := s.r.Get(ctx, uid)
	s.c.Set(order.OrderUID, order, s.ttlCache)
	if err != nil {
		log.Printf("[err] service:%s\n", err.Error())
		return nil, err
	}

	return order, nil
}

func (s *service) StoreOrder(ctx context.Context, order *domain.Order) (string, error) {
	uid, err := s.r.Store(ctx, order)
	if err != nil {
		log.Printf("[err] service:%s\n", err.Error())
		return "", err
	}
	if err == nil {
		s.c.Set(order.OrderUID, order, s.ttlCache)
	}

	return uid, nil
}
