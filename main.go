package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// criando um novo roteador
	router := chi.NewRouter()
	// configurando o middleware p/ logar as requisições
	router.Use(middleware.Logger)

	// criando uma rota simples
	router.Get("/ola", basicHandler)

	// criando um servidor
	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	// error handling
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Falha crítica!", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Olá!"))
}
