package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

// cria uma nova instância de App com um roteador carregado
func New() *App {
	app := &App{
		router: loadRouter(), // loadRouter, from routes.go
		rdb: redis.NewClient(&redis.Options{}),
	}

	return app
}

// inicia o servidor
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// err handling
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("falha na conexão com o redis: %w", err)
	}

	fmt.Println("Iniciando o servidor...")

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor: %w", err)
	}

	return nil
}
