package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Afonso247/api-pedidos/model"
	"github.com/Afonso247/api-pedidos/repositorio/pedido"
)

type Pedido struct{
	Repo *pedido.RedisRepo
}

// funções que respondem a requests HTTP

// POST request
func (p *Pedido) Create(w http.ResponseWriter, r *http.Request) {
	// obtem o body da requisição
	var body struct {
		ClienteID uuid.UUID `json:"cliente_id"`
		LineItems []model.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Define o fuso horário UTC-3
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Erro ao carregar a localização:", err)
		return
	}

	now := time.Now().In(loc)

	// cria o pedido
	pedido := model.Pedido{
		PedidoID: rand.Uint64(),
		ClienteID: body.ClienteID,
		LineItems: body.LineItems,
		CriadoEm: &now,
	}

	// insere o pedido
	err = p.Repo.Insert(r.Context(), &pedido)
	if err != nil {
		fmt.Println("Erro ao inserir o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// se o pedido for inserido com sucesso, retorna o pedido
	res, err := json.Marshal(pedido)
	if err != nil {
		fmt.Println("Erro ao serializar o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

// GET request
func (p *Pedido) List(w http.ResponseWriter, r *http.Request) {
	// obtem o cursor da requisição (default = 0)
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	// converte cursor em uint64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// busca os pedidos com base no cursor e limit
	const limit = 50
	res, err := p.Repo.FindAll(r.Context(), pedido.FindAllPage{
		Offset: cursor,
		Limit:  limit,
	})
	if err != nil {
		fmt.Println("Erro ao buscar os pedidos:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// cria response p/ armazenar os pedidos encontrados
	var response struct {
		Items []model.Pedido `json:"pedidos"`
		Prox  uint64         `json:"prox,omitempty"`
	}

	response.Items = res.Pedidos
	response.Prox = res.Cursor

	// serializa os pedidos
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Erro ao serializar os pedidos:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// GET request
func (p *Pedido) GetByID(w http.ResponseWriter, r *http.Request) {
	// obtem o ID da requisição
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	// converte idParam em uint64
	id, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// busca o pedido pelo ID
	o, err := p.Repo.FindByID(r.Context(), id)
	if errors.Is(err, pedido.ErrNotExist) {
		// se o pedido não existir, retorna 404
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		// se houver algum outro erro, retorna 500
		fmt.Println("Erro ao buscar o pedido pelo id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// serializa o pedido
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("Erro ao serializar o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// PUT request
func (p *Pedido) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// obtem o ID da requisição
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	// converte idParam em uint64
	id, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// busca o pedido pelo ID
	elPedido, err := p.Repo.FindByID(r.Context(), id)
	if errors.Is(err, pedido.ErrNotExist) {
		// se o pedido não existir, retorna 404
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		// se houver algum outro erro, retorna 500
		fmt.Println("Erro ao buscar o pedido pelo id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const enviadoStatus = "enviado"
	const concluidoStatus = "concluido"

	// Define o fuso horário UTC-3
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Erro ao carregar a localização:", err)
		return
	}

	now := time.Now().In(loc)

	// atualiza o pedido com base no request body (status)
	switch body.Status {
	case enviadoStatus:
		// se o status for "enviado" e o pedido não tiver sido enviado, 
		// atualiza o campo EnviadoEm
		if elPedido.EnviadoEm != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		elPedido.EnviadoEm = &now
	case concluidoStatus:
		// se o status for "concluido" e o pedido não tiver sido enviado ou 
		// não tiver sido concluído, atualiza o campo ConcluidoEm
		if elPedido.ConcluidoEm != nil || elPedido.EnviadoEm == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		elPedido.ConcluidoEm = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// atualiza o pedido
	err = p.Repo.UpdateByID(r.Context(), &elPedido)
	if err != nil {
		fmt.Println("Erro ao atualizar o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// serializa o pedido
	if err := json.NewEncoder(w).Encode(elPedido); err != nil {
		fmt.Println("Erro ao serializar o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// DELETE request
func (p *Pedido) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// obtem o ID da requisição
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	// converte idParam em uint64
	id, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// deleta o pedido pelo ID
	err = p.Repo.DeleteByID(r.Context(), id)
	if errors.Is(err, pedido.ErrNotExist) {
		// se o pedido não existir, retorna 404
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		// se houver algum outro erro, retorna 500
		fmt.Println("Erro ao deletar o pedido:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
