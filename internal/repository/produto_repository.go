package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"
	models "github.com/julio-pupim/lojaestoque/internal/domain"
)

// ProdutoRepository encapsula acessos ao banco para produtos
type ProdutoRepository struct {
	db *sql.DB
}

// NewProdutoRepository cria uma instância de ProdutoRepository
func NewProdutoRepository(db *sql.DB) *ProdutoRepository {
	return &ProdutoRepository{db: db}
}

// Save insere ou atualiza (soma estoque) de um produto
func (r *ProdutoRepository) Save(p *models.Produto) error {
	// 1) Busca pelo par (fornecedor_id, codigo_fornecedor)
	var id, qtdAtual int64
	row := r.db.QueryRow(
		`SELECT id, quantidade_estoque
           FROM produtos
          WHERE fornecedor_id = ? AND codigo_fornecedor = ?`,
		p.Fornecedor.Id, p.CodigoFornecedor,
	)
	err := row.Scan(&id, &qtdAtual)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// 2a) Não existe → inserção
		_, err = r.db.Exec(
			`INSERT INTO produtos
               (nome, fornecedor_id, codigo_fornecedor,
                quantidade_estoque, preco)
             VALUES (?, ?, ?, ?, ?)`,
			p.Nome,
			p.Fornecedor.Id,
			p.CodigoFornecedor,
			p.QuantidadeEstoque,
			p.Preco.String(),
		)

	case err != nil:
		// 2b) Erro inesperado no SELECT
		return fmt.Errorf("erro ao buscar produto existente: %w", err)

	default:
		// 2c) Já existe → soma estoque e atualiza preço
		novaQtde := qtdAtual + p.QuantidadeEstoque
		_, err = r.db.Exec(
			`UPDATE produtos
                SET quantidade_estoque = ?, preco = ?
              WHERE id = ?`,
			novaQtde,
			p.Preco.String(),
			id,
		)
	}

	return err
}

// Find consulta produtos com filtros opcionais
func (r *ProdutoRepository) Find(filters map[string]any, limit, offset int) ([]models.Produto, error) {
	base := `SELECT id, nome, fornecedor_id, codigo_fornecedor, quantidade_estoque, preco FROM produtos`
	var clauses []string
	var args []any

	if v, ok := filters["nome"]; ok {
		clauses = append(clauses, "nome LIKE ?")
		args = append(args, "%"+v.(string)+"%")
	}
	if v, ok := filters["fornecedor_id"]; ok {
		clauses = append(clauses, "fornecedor_id = ?")
		args = append(args, v)
	}
	if v, ok := filters["codigo_fornecedor"]; ok {
		clauses = append(clauses, "codigo_fornecedor = ?")
		args = append(args, v)
	}
	if v, ok := filters["preco_min"]; ok {
		clauses = append(clauses, "preco >= ?")
		args = append(args, v.(*apd.Decimal).String())
	}
	if v, ok := filters["preco_max"]; ok {
		clauses = append(clauses, "preco <= ?")
		args = append(args, v.(*apd.Decimal).String())
	}

	if len(clauses) > 0 {
		base += " WHERE " + strings.Join(clauses, " AND ")
	}
	base += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var (
			id, fornecedorID, qtdEstoque, qtdMinima int64
			nome, codForn, codBarras, precoStr      string
		)
		if err := rows.Scan(
			&id, &nome, &fornecedorID, &codForn,
			&codBarras, &qtdEstoque, &precoStr, &qtdMinima,
		); err != nil {
			return nil, err
		}
		forn := models.Fornecedor{Id: fornecedorID, Nome: nome}
		p, err := models.NewProduto(
			id,
			nome,
			forn,
			codForn,
			qtdEstoque,
			precoStr,
		)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, *p)
	}
	return produtos, rows.Err()
}

// Delete remove um produto pelo ID
func (r *ProdutoRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM produtos WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Update edita campos (exceto ID, histórico de estoque)
func (r *ProdutoRepository) Update(p *models.Produto) error {
	// assume que Validate foi chamada antes
	_, err := r.db.Exec(
		`UPDATE produtos SET nome = ?, codigo_fornecedor = ?, quantidade_estoque = ?, preco = ?
		 WHERE fornecedor_id = ? AND codigo_fornecedor = ?`,
		p.Nome, p.CodigoFornecedor, p.QuantidadeEstoque, p.Preco, p.Fornecedor.Id, p.CodigoFornecedor,
	)
	return err
}

// --- Handler de produtos (pseudocódigo, importar chi, handler, etc.) ---
// r.Route("/produtos", func(r chi.Router) {
//   pr := repository.NewProdutoRepository(db)
//   r.Post("/", handler.CreateOrAddProduto(pr))
//   r.Get("/", handler.SearchProdutos(pr))
//   r.Delete("/{id}", handler.DeleteProduto(pr))
//   r.Patch("/{id}", handler.UpdateProduto(pr))
// })

// Os handler functions gerariam closures recebendo pr *ProdutoRepository e retornando http.HandlerFunc
