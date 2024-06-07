package service

import (
	"context"
	"encoding/json"

	"github.com/Artenso/wb-l0/internal/cache"
	"github.com/Artenso/wb-l0/internal/model"
	"github.com/Artenso/wb-l0/internal/repository/postgres"
)

// IService working with service
type IService interface {
	AddOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
	RestoreCache(ctx context.Context) error
}

type service struct {
	repository postgres.IRepository
	cache      cache.ICache
}

// New createse new service
func New(repository postgres.IRepository, cache cache.ICache) IService {
	return &service{
		repository: repository,
		cache:      cache,
	}

}

// AddOrder adds order to db and cache
func (s *service) AddOrder(ctx context.Context, order *model.Order) error {
	if err := s.repository.AddOrder(ctx, order); err != nil {
		return err
	}
	if err := s.cache.AddOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

// GetOrder gets order from cache
func (s *service) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	order, err := s.cache.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// RestoreCache load orders from db to cache
func (s *service) RestoreCache(ctx context.Context) error {
	orders, err := s.repository.GetOrders(ctx)
	if err != nil {
		return err
	}
	var jsonOrder []byte
	var order *model.Order

	for orders.Next() {
		err := orders.Scan(&jsonOrder)
		if err != nil {
			return err
		}
		json.Unmarshal(jsonOrder, &order)
		err = s.cache.AddOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
