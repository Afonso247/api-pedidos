package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Afonso247/api-pedidos/handler"
)

// carrega o roteador
func loadRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	// HTTP Status 200
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// loadPedidoRouter lida com os requests HTTP na subrota /pedidos
	router.Route("/pedidos", loadPedidoRouter)

	return router
}

func loadPedidoRouter(router chi.Router) {
	// instancia do package handler
	pedidoHandler := &handler.Pedido{}

	// HTTP requests
	router.Post("/", pedidoHandler.Create)
	router.Get("/", pedidoHandler.List)
	router.Get("/{id}", pedidoHandler.GetByID)
	router.Put("/{id}", pedidoHandler.UpdateByID)
	router.Delete("/{id}", pedidoHandler.DeleteByID)
}
