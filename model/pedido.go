package model

import "github.com/google/uuid"

// modelo do pedido
type Pedido struct{

	PedidoID uint64 `json:"pedido_id"`
	ClienteID uuid.UUID `json:"cliente_id"`
	LineItems []LineItem `json:"line_items"`
	CriadoEm *time.Time `json:"criado_em"`
	EnviadoEm *time.Time `json:"enviado_em"`
	ConcluidoEm *time.Time `json:"concluido_em"`
}

// modelo da linha do pedido
type LineItem struct{
	
	ItemID uuid.UUID `json:"item_id"`
	Quantidade uint `json:"quantidade"`
	Preco uint `json:"preco"`
}