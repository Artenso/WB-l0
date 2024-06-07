package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Artenso/wb-l0/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/stan.go"
)

// IProducer working with producer
type IProducer interface {
	Run() error
	Stop()
}

type producer struct {
	connection stan.Conn
}

// New creates new producer
func New() IProducer {
	// connect to nats-streaming
	conn, err := stan.Connect("test-cluster", "producer")
	if err != nil {
		log.Fatalf("producer failed to create nats-streaming connection: %s", err.Error())
	}

	return &producer{
		connection: conn,
	}
}

// Run starts producer
func (p *producer) Run() error {
	order := new(model.Order)

	for {
		time.Sleep(time.Second * 2)

		// create new order
		gofakeit.Struct(order)

		// convert model to json
		jsonOrder, err := json.MarshalIndent(order, "", " ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %s", err.Error())
		}

		// sand json to nats
		err = p.connection.Publish("JsonPipe", jsonOrder)
		if err != nil {
			return fmt.Errorf("failed to send order: %s", err.Error())
		} else {
			log.Printf("Send successfully: %s", order.Order_uid)
		}
	}
}

// Stop stops producer and close connection
func (p *producer) Stop() {
	p.connection.Close()
}
