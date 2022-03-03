package service

import (
	"context"
	"github.com/wubba-com/L0/internal/app/domain"
	"github.com/wubba-com/L0/pkg/cache"
	"strconv"
	"time"
)

func NewItemService(repository domain.ItemRepository, cache cache.Cache, ttl time.Duration) domain.ItemService {
	return &serviceItem{repository, cache, ttl}
}

type serviceItem struct {
	r        domain.ItemRepository
	c        cache.Cache
	ttlCache time.Duration
}

func (s *serviceItem) GetItemByOrderUID(ctx context.Context, orderUID string) ([]*domain.Item, error) {
	items, err := s.r.GetByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *serviceItem) StoreItem(ctx context.Context, item *domain.Item) (uint64, error) {
	uid, err := s.r.Store(ctx, item)
	if err != nil {
		return 0, err
	}

	uidCache := strconv.FormatUint(uid, 10)
	s.c.Set(uidCache, item, s.ttlCache)

	return uid, err
}
