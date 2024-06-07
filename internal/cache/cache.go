package cache

import (
	"context"
	"fmt"

	"github.com/Artenso/wb-l0/internal/model"
)

// ICache working with cache
type ICache interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
}

type cache struct {
	storage map[string]model.Order
}

// New creates new cache
func New() ICache {
	return &cache{
		storage: make(map[string]model.Order),
	}
}

// AddOrder adds order to cache
func (c *cache) AddOrder(ctx context.Context, order *model.Order) error {
	c.storage[order.Order_uid] = *order
	return nil
}

// GetOrder gets order from cache
func (c *cache) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	if order, ok := c.storage[orderUID]; ok {
		return &order, nil
	}

	return nil, fmt.Errorf("order not found")
}
