package handler

import (
	"fmt"
	"net/http"
)

type Pedido struct{}

// HTTP requests
func (p *Pedido) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pedido criado com sucesso!")
}

func (p *Pedido) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listando todos os pedidos...")
}

func (p *Pedido) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pegando um pedido pelo ID...")
}

func (p *Pedido) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Atualizando um pedido pelo ID...")
}

func (p *Pedido) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deletando um pedido pelo ID...")
}
