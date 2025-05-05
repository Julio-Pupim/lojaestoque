package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/julio-pupim/lojaestoque/internal/domain"
)

type ClienteRepository struct {
	db *sql.DB
}

func NewClienteRepository(db *sql.DB) *ClienteRepository {
	return &ClienteRepository{db: db}
}

func (cr *ClienteRepository) BuscarClientes(filters map[string]any, limit, offset int) ([]domain.Cliente, error) {
	sql := "SELECT * FROM clientes "
	var clauses []string
	var args []any

	if v, ok := filters["nome"]; ok {
		clauses = append(clauses, "nome ILIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if v, ok := filters["telefone"]; ok {
		clauses = append(clauses, "telefone LIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if len(clauses) > 0 {
		sql += " WHERE " + strings.Join(clauses, " AND ")
	}
	sql += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := cr.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []domain.Cliente
	for rows.Next() {
		var (
			id                           int64
			nome, telefone, dataCadastro string
		)
		if err := rows.Scan(&id, &nome, &telefone, &dataCadastro); err != nil {
			fmt.Print(err)
			return nil, err
		}
		clientes = append(clientes, domain.Cliente{ID: id, Nome: nome, Telefone: telefone, DataCadastro: dataCadastro})
	}

	return clientes, nil
}

func (cr *ClienteRepository) SalvarCliente(c *domain.Cliente) (sql.Result, error) {
	result, err := cr.db.Exec("INSERT INTO clientes (nome, telefone, data_cadastro) VALUES (?,?,?)",
		c.Nome, c.Telefone, c.DataCadastro)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cr *ClienteRepository) DeletarCliente(id int64) error {
	_, err := cr.db.Exec("DELETE FROM clientes WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ClienteRepository) BuscarClientePorId(id int64) (domain.Cliente, error) {
	row := cr.db.QueryRow("SELECT * FROM clientes WHERE id = ? ", id)
	if row.Err() != nil {
		return domain.Cliente{}, row.Err()
	}
	var (
		nome, telefone, dataCadastro string
	)
	err := row.Scan(nome, telefone, dataCadastro)
	if err != nil {
		return domain.Cliente{}, err
	}
	return domain.Cliente{ID: id, Nome: nome, Telefone: telefone, DataCadastro: dataCadastro}, nil
}

func (cr *ClienteRepository) AtualizarCliente(id int64, cliente domain.Cliente) (sql.Result, error) {
	var clauses []string

	if v := cliente.Nome; v != "" {
		clauses = append(clauses, "nome = ?")
	}
	if v := cliente.Telefone; v != "" {
		clauses = append(clauses, "telefone = ?")
	}
	sql := "UPDATE clientes SET "
	sql += strings.Join(clauses, ", ")
	sql += "WHERE id = ? "
	result, err := cr.db.Exec(sql, cliente.Nome, cliente.Telefone, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
