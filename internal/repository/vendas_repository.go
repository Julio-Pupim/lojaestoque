package repository

import (
	"database/sql"

	"github.com/julio-pupim/lojaestoque/internal/domain"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}
func (r *SaleRepository) Save(sale *domain.Sale) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(
		`INSERT INTO sales (cliente_id, data_venda, total, status_pagamento) VALUES (?, ?, ?, ?)`,
		sale.ClientID, sale.DataVenda, sale.Total, sale.PaymentStatus,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	saleID, _ := res.LastInsertId()
	sale.ID = saleID

	for i := range sale.Items {
		item := &sale.Items[i]
		res, err := tx.Exec(
			`INSERT INTO sale_items (venda_id, produto_id, quantidade, preco_unitario, total) VALUES (?, ?, ?, ?, ?)`,
			saleID, item.ProductID, item.Quantity, item.UnitPrice, item.Total,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		itemID, _ := res.LastInsertId()
		item.ID = itemID
		item.SaleID = saleID
	}

	return tx.Commit()
}
