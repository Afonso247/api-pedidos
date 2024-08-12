package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Afonso247/api-pedidos/app"
)

// inicia a aplicação a partir do package app
func main() {
	app := app.New()

	// listen por interrupções (ex: SIGINT, SIGTERM)
	// ao receber, cancele o contexto
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("falha em iniciar a aplicação:", err)
	}
}
