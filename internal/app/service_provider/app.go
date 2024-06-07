package serviceprovider

import (
	"context"
	"net/http"
	"sync"

	"github.com/Artenso/wb-l0/internal/config"
)

const (
	pathToConfig = "config.json"
)

// App is the application struct
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

// NewApp creates new App
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run starts app
func (a *App) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}

	if err := a.serviceProvider.getService(ctx).RestoreCache(ctx); err != nil {
		return err
	}

	if err := a.serviceProvider.getConsumer(ctx).Subscribe(ctx); err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.httpServer.ListenAndServe()
	}()

	wg.Wait()
	return nil
}

// Stop shutdown all connections and services
func (a *App) Stop(ctx context.Context) {
	a.serviceProvider.consumer.Stop(ctx)
	a.serviceProvider.nsConn.Close()
	a.httpServer.Shutdown(ctx)
}

// initDeps initialize dependencies
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// initConfig loads service configuration
func (a *App) initConfig(_ context.Context) error {
	config.Read(pathToConfig)
	return nil
}

// initServiceProvider initialize serviceProvider
func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

// initHttpServer initialize http server
func (a *App) initHttpServer(ctx context.Context) error {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/order/", a.serviceProvider.getHandler(ctx).GetOrder)
	httpMux.Handle("/", http.FileServer(http.Dir("./static/ui.html")))
	a.httpServer = &http.Server{
		Addr:    config.GetHttpPort(),
		Handler: httpMux,
	}

	return nil
}
