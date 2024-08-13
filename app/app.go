package app

import (
	"context"
	"fmt"
	"time"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	config Config
}

// cria uma nova instância de App com um roteador carregado
func New(config Config) *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddr,
		}),
		config: config,
	}
	
	app.loadRouter()

	return app
}

// inicia o servidor
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	// err handling
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("falha na conexão com o redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("falha ao fechar o redis:", err)
		}
	}()

	fmt.Println("Iniciando o servidor...")

	ch := make(chan error, 1)

	// ch <- err
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("falha ao iniciar o servidor: %w", err)
		}
		close(ch)
	}()

	// aguarda pelo contexto
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}
