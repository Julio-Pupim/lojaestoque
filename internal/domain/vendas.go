package domain

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDENTE"
	PaymentStatusPaid    PaymentStatus = "PAGO"
	PaymentStatusPartial PaymentStatus = "PARCIAL"
)

type Sale struct {
	ID            int64         `json:"id"`
	ClientID      int64         `json:"cliente_id"`
	DataVenda     time.Time     `json:"data"`
	Total         Decimal       `json:"total"`
	PaymentDate   *time.Time    `json:"data_pagamento,omitempty"`
	PaymentStatus PaymentStatus `json:"status_pagamento"`
	Items         []SaleItem    `json:"items"`
}

// SaleItem representa um item de uma venda
type SaleItem struct {
	ID        int64   `json:"id"`
	SaleID    int64   `json:"venda_id"`
	ProductID int64   `json:"produto_id"`
	Quantity  int     `json:"quantidade"`
	UnitPrice Decimal `json:"preco_unitario"`
	Total     Decimal `json:"total"`
}
