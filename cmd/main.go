package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	serviceProvider "github.com/Artenso/wb-l0/internal/app/service_provider"
	producer "github.com/Artenso/wb-l0/pkg/producer/nats"
)

func main() {
	ctx := context.Background()
	wg := sync.WaitGroup{}

	// create producer
	producer := producer.New()

	// create application
	app, err := serviceProvider.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	wg.Add(1)
	go func() {
		// Make a signal channel. Register SIGINT.
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt)

		// Wait for the signal.
		<-sigch
		// stop service
		app.Stop(ctx)
		producer.Stop()
	}()

	// run producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := producer.Run(); err != nil {
			log.Fatalf("failed to run producer: %s", err.Error())
		}
	}()

	// run application
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = app.Run(ctx); err != nil {
			log.Fatalf("failed to run app: %s", err.Error())
		}
	}()

	wg.Wait()
}
