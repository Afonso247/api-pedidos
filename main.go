package main

import (
	"context"
	"fmt"

	"github.com/Afonso247/api-pedidos/app"
)

// inicia a aplicação a partir do package app
func main() {
	app := app.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("falha em iniciar a aplicação:", err)
	}
}
