package service

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/cache"
	"time"
)

func NewPaymentService(repository domain.PaymentRepository, cache cache.Cache, ttl time.Duration) domain.PaymentService {
	return &servicePayment{repository, cache, ttl}
}

type servicePayment struct {
	r domain.PaymentRepository
	c cache.Cache
	ttlCache time.Duration
}

func (s *servicePayment) StorePayment(ctx context.Context, payment *domain.Payment) (string, error) {
	uid, err := s.r.Store(ctx, payment)
	if err != nil {
		return "", err
	}

	s.c.Set(uid, payment, s.ttlCache)

	return uid, err
}

