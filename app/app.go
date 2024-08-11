package app

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

// cria uma nova instância de App com um roteador carregado
func New() *App {
	app := &App{
		router: loadRouter(), // loadRouter, from routes.go
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
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("falha ao iniciar o servidor: %w", err)
	}

	return nil
}
