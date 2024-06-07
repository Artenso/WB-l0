package serviceprovider

import (
	"context"
	"log"

	"github.com/Artenso/wb-l0/internal/cache"
	"github.com/Artenso/wb-l0/internal/config"
	"github.com/Artenso/wb-l0/internal/consumer/nats"
	"github.com/Artenso/wb-l0/internal/http/handlers"
	"github.com/Artenso/wb-l0/internal/repository/postgres"
	"github.com/Artenso/wb-l0/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
)

// serviceProvider di-container
type serviceProvider struct {
	dbConn     *pgxpool.Pool
	nsConn     stan.Conn
	repository postgres.IRepository
	cache      cache.ICache
	service    service.IService
	consumer   nats.IConsumer
	handler    handlers.IHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// getDbConn creates connection with database
func (s *serviceProvider) getDbConn(ctx context.Context) *pgxpool.Pool {
	if s.dbConn == nil {
		conn, err := pgxpool.New(ctx, config.GetDbDSN())
		if err != nil {
			log.Fatalf("failed to init db connection: %s", err.Error())
		}

		s.dbConn = conn
	}

	return s.dbConn
}

// getNsConn creates connection with nats-streaming
func (s *serviceProvider) getNsConn(_ context.Context) stan.Conn {
	if s.nsConn == nil {
		conn, err := stan.Connect(config.GetClusterID(), config.GetClientID())
		if err != nil {
			log.Fatalf("failed to init nats-streaming connection: %s", err.Error())
		}

		s.nsConn = conn
	}

	return s.nsConn
}

func (s *serviceProvider) getRepository(ctx context.Context) postgres.IRepository {
	if s.repository == nil {
		s.repository = postgres.New(s.getDbConn(ctx))
	}

	return s.repository
}

func (s *serviceProvider) getCache() cache.ICache {
	if s.cache == nil {
		s.cache = cache.New()
	}

	return s.cache
}

func (s *serviceProvider) getService(ctx context.Context) service.IService {
	if s.service == nil {
		s.service = service.New(s.getRepository(ctx), s.getCache())
	}

	return s.service
}

func (s *serviceProvider) getConsumer(ctx context.Context) nats.IConsumer {
	if s.consumer == nil {
		s.consumer = nats.New(s.getNsConn(ctx), s.getService(ctx))
	}

	return s.consumer
}

func (s *serviceProvider) getHandler(ctx context.Context) handlers.IHandler {
	if s.handler == nil {
		s.handler = handlers.New(s.getService(ctx))
	}

	return s.handler
}
