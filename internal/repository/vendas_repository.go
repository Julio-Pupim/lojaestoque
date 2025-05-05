package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/julio-pupim/lojaestoque/internal/domain"
)

type VendasRepository struct {
	db *sql.DB
}

func NewVendasRepository(db *sql.DB) *VendasRepository {
	return &VendasRepository{db: db}
}
func (vr *VendasRepository) SalvarVenda(sale *domain.Sale) error {
	tx, err := vr.db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(
		`INSERT INTO vendas (cliente_id, data_venda, total, status_pagamento) VALUES (?, ?, ?, ?)`,
		sale.ClientID, sale.DataVenda, sale.Total, sale.PaymentStatus,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	saleID, _ := res.LastInsertId()
	sale.ID = saleID
	valueStrings := make([]string, 0, len(sale.Items))
	valueArgs := make([]any, 0, len(sale.Items)*5)

	for _, item := range sale.Items {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			saleID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.Total)
	}
	query := fmt.Sprintf(`
  INSERT INTO vendas_produtos
    (venda_id, produto_id, quantidade, preco_unitario, total)
  VALUES %s`, strings.Join(valueStrings, ","))
	if _, err := tx.Exec(query, valueArgs...); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (vr *VendasRepository) BuscarVendas(filters map[string]any, limit, offset int) ([]domain.Sale, error) {
	sql := `SELECT v.* FROM vendas v 
	join clientes c on c.id = v.cliente_id 
	join vendas_produtos vp on vp.venda_id = v.id
	join produto p on vp.produto_id = p.id 
	`
	var clauses []string
	var args []any
	if v, ok := filters["nome_cliente"]; ok {
		clauses = append(clauses, "c.nome ILIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if v, ok := filters["nome_produto"]; ok {
		clauses = append(clauses, "p.nome ILIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if v, ok := filters["total"]; ok {
		clauses = append(clauses, "v.total = ?")
		args = append(args, v.(string))
	}
	if v, ok := filters["data_pagamento"]; ok {
		clauses = append(clauses, "v.data_pagamento = ?")
		args = append(args, v.(string))
	}

	if v, ok := filters["data_venda"]; ok {
		clauses = append(clauses, "v.data_venda = ?")
		args = append(args, v.(string))
	}
	if v, ok := filters["status_pagamento"]; ok {
		clauses = append(clauses, "v.status_pagamento ILIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if len(clauses) > 0 {
		sql += " WHERE " + strings.Join(clauses, " AND ")
	}
	sql += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := vr.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vendas []domain.Sale
	for rows.Next() {
		var (
			id, clienteId            int64
			dataVenda, dataPagamento time.Time
			total                    domain.Decimal
			statusPagamento          domain.PaymentStatus
		)
		if err := rows.Scan(id, clienteId, dataVenda, total, dataPagamento, statusPagamento); err != nil {
			return nil, err
		}
		vendas = append(vendas, domain.Sale{})
	}
	return vendas, nil
}
func (vr *VendasRepository) buscarVendaPorId(id int64) (domain.Sale, error) {
	row := vr.db.QueryRow("SELECT * FROM vendas WHERE id = ?", id)
	if row.Err() != nil {
		return domain.Sale{}, row.Err()
	}
	var (
		idResponse, clienteId int64
		dataVenda             time.Time
		dataPagamento         *time.Time
		total                 domain.Decimal
		status_pagamento      domain.PaymentStatus
	)
	err := row.Scan(idResponse, clienteId, dataVenda, dataPagamento, total, status_pagamento)
	if err != nil {
		return domain.Sale{}, err
	}

	return domain.Sale{ID: idResponse, ClientID: clienteId, DataVenda: dataVenda, Total: total,
		PaymentDate: dataPagamento, PaymentStatus: status_pagamento}, nil
}

// func (vr *VendasRepository) updateVenda() (domain.Sale, error) {
// }

func (vr *VendasRepository) DeletarVenda(id int64) error {
	_, err := vr.db.Exec("DELETE FROM vendas WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
