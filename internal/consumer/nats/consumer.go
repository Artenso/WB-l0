package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Artenso/wb-l0/internal/model"
	"github.com/Artenso/wb-l0/internal/service"
	"github.com/nats-io/stan.go"
)

// IConsumer working with consumer
type IConsumer interface {
	Subscribe(ctx context.Context) error
	Stop(ctx context.Context) error
}

type consumer struct {
	connection   stan.Conn
	subscription stan.Subscription
	service      service.IService
}

// New creates new consumer
func New(conn stan.Conn, service service.IService) IConsumer {
	return &consumer{
		connection: conn,
		service:    service,
	}

}

// Subscribe create subscription to nats-streaming
func (c *consumer) Subscribe(ctx context.Context) error {
	order := new(model.Order)

	subscription, err := c.connection.Subscribe("JsonPipe", func(m *stan.Msg) {
		err := json.Unmarshal(m.Data, order)
		if err != nil {
			log.Println(err)
		} else {
			err = c.service.AddOrder(ctx, order)
			if err != nil {
				log.Println(err)
			}
		}
	}, stan.DeliverAllAvailable())
	if err != nil {
		return err
	}
	c.subscription = subscription
	return nil
}

// Stop drops subscription to nats-streaming
func (c *consumer) Stop(ctx context.Context) error {
	if err := c.subscription.Unsubscribe(); err != nil {
		return err
	}

	return nil
}
