package pedido

import (
	"context"
	"fmt"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/Afonso247/api-pedidos/model"
)

// RedisRepo implementa o repositorio de pedidos
type RedisRepo struct {
	Client *redis.Client
}

// gera uma chave unica para um pedido
func pedidoIDKey(id uint64) string {
	return fmt.Sprintf("pedido:%d", id)
}

// insere um pedido
func (r *RedisRepo) Insert(ctx context.Context, pedido *model.Pedido) error {
	// serializa o pedido
	data, err := json.Marshal(pedido)
	if err != nil {
		return fmt.Errorf("falha ao serializar o pedido: %w", err)
	}

	// pega uma chave unica pro pedido
	key := pedidoIDKey(pedido.PedidoID)

	// cria um pipeline de transação
	// é responsável para executar múltiplos comandos de forma atômica
	txn := r.Client.TxPipeline()

	// seta o pedido, mas somente se o pedido ainda não existir
	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("falha ao gravar o pedido: %w", err)
	}

	// adiciona o pedido a lista de pedidos
	if err := txn.SAdd(ctx, "pedidos", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("falha ao adicionar o pedido à lista: %w", err)
	}

	// executa o pipeline, e retorna qualquer erro que ocorrer durante a execução
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("falha de execução do pipeline: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("o pedido não existe")

// procura um pedido pelo ID
func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Pedido, error) {
	key := pedidoIDKey(id)

	// procura o pedido
	val, err := r.Client.Get(ctx, key).Result()

	// se o pedido não for encontrado, retorna ErrNotExist
	if errors.Is(err, redis.Nil) {
		return model.Pedido{}, ErrNotExist
	// se outro erro ocorrer, retorna o erro
	} else if err != nil {
		return model.Pedido{}, fmt.Errorf("falha ao buscar o pedido: %w", err)
	}

	// se o pedido for encontrado, desserializa o pedido para um objeto model.Pedido
	var pedido model.Pedido
	err = json.Unmarshal([]byte(val), &pedido)
	// se ocorrer um erro, retorna o erro
	if err != nil {
		return model.Pedido{}, fmt.Errorf("falha ao desserializar o pedido: %w", err)
	}
	// note que em todos os casos de erro, é retornado o pedido vazio

	return pedido, nil
}

// deleta um pedido pelo ID
func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := pedidoIDKey(id)

	// cria um pipeline de transação
	// é responsável para executar múltiplos comandos de forma atômica
	txn := r.Client.TxPipeline()

	// deleta o pedido
	err := txn.Del(ctx, key).Err()

	// se o pedido não existir, retorna ErrNotExist
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	// se ocorrer outro erro, retorna o erro
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("falha ao deletar o pedido: %w", err)
	}

	// remove o pedido da lista de pedidos
	if err := txn.SRem(ctx, "pedidos", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("falha ao remover o pedido da lista: %w", err)
	}

	// executa o pipeline, e retorna qualquer erro que ocorrer durante a execução
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("falha de execução do pipeline: %w", err)
	}

	return nil
}

// atualiza um pedido pelo ID
func (r *RedisRepo) UpdateByID(ctx context.Context, pedido *model.Pedido) error {
	// serializa o pedido
	data, err := json.Marshal(pedido)
	if err != nil {
		return fmt.Errorf("falha ao serializar o pedido: %w", err)
	}

	// pega uma chave unica pro pedido
	key := pedidoIDKey(pedido.PedidoID)

	// atualiza o pedido somente se o pedido existir
	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("falha ao atualizar o pedido: %w", err)
	}

	return nil
}

// struct de paginação
type FindAllPage struct {
	Offset int
	Limit  int
}

// struct de resultado p/ FindAll
type FindResult struct {
	Pedidos []model.Pedido
	Cursor   uint64
}

// procura todos os pedidos
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) ([]model.Pedido, error) {

	// procura os pedidos pelo Offset e Limit
	res := r.Client.SScan(ctx, "pedidos", page.Offset, "*", int64(page.Limit))

	// retorna um cursor e algumas chaves
	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("falha ao buscar os pedidos pelo id: %w", err)
	}

	// se nenhuma chave for encontrada, retorna um FindResult vazio
	if len(keys) == 0 {
		return FindResult{
			Pedidos: []model.Pedido{},
			}, nil
	}

	// usa MGet para obter os pedidos a partir das chaves
	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("falha ao buscar os pedidos: %w", err)
	}

	pedidos := make([]model.Pedido, len(xs))

	// deserializa cada pedido para um objeto de model.Pedido
	for i, x := range xs {
		x := x.(string)
		var pedido model.Pedido
		
		err = json.Unmarshal([]byte(x), &pedido)
		if err != nil {
			return FindResult{}, fmt.Errorf("falha ao desserializar o pedido em json: %w", err)
		}

		pedidos[i] = pedido
	}

	// retorna FindResult com os pedidos encontrados e o cursor
	return FindResult{
		Pedidos: pedidos, 
		Cursor: cursor,
		}, nil
}