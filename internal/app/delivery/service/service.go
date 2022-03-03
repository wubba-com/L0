package service

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/cache"
	"strconv"
	"time"
)

func NewDeliveryService(repository domain.DeliveryRepository, cache cache.Cache, ttl time.Duration) domain.DeliveryService {
	return &serviceDelivery{repository, cache, ttl}
}

type serviceDelivery struct {
	r domain.DeliveryRepository
	c cache.Cache
	ttlCache time.Duration
}

func (s *serviceDelivery) GetDeliveryByOrderUID(ctx context.Context, orderUID string) (*domain.Delivery, error)  {
	delivery, err := s.r.GetByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	uidCache := strconv.FormatUint(delivery.DeliveryUID, 10)
	s.c.Set(uidCache, delivery, s.ttlCache)

	return delivery, nil
}

func (s *serviceDelivery) StoreDelivery(ctx context.Context, delivery *domain.Delivery) (uint64, error) {
	uid, err := s.r.Store(ctx, delivery)
	if err != nil {
		return 0,  err
	}

	uidCache := strconv.FormatUint(uid, 10)
	s.c.Set(uidCache, delivery, s.ttlCache)

	return uid, nil
}
