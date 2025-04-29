package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cockroachdb/apd/v3"
	"github.com/go-chi/chi/v5"
	"github.com/julio-pupim/lojaestoque/internal/domain"
	"github.com/julio-pupim/lojaestoque/internal/middleware"
	"github.com/julio-pupim/lojaestoque/internal/repository"
)

// CreateOrAddProduto retorna um http.HandlerFunc que faz INSERT ou atualiza estoque
func CreateOrAddProduto(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto struct {
			Nome              string `json:"nome"`
			FornecedorID      int64  `json:"fornecedor_id"`
			CodigoFornecedor  string `json:"codigo_fornecedor"`
			QuantidadeEstoque int64  `json:"quantidade_estoque"`
			Preco             string `json:"preco"`
			QtdMinima         int64  `json:"qtd_minima"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Constrói domain model
		forn := domain.Fornecedor{Id: dto.FornecedorID}
		p, err := domain.NewProduto(
			0,
			dto.Nome,
			forn,
			dto.CodigoFornecedor,
			dto.QuantidadeEstoque,
			dto.Preco,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Persiste
		if err := pr.Save(p); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao salvar produto: %v", err), http.StatusInternalServerError)
			return
		}

		// Retorna criado/atualizado
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	}
}

// SearchProdutos retorna http.HandlerFunc que busca com filtros e paginação
func SearchProdutos(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai filtros
		filters := make(map[string]any)
		if v := r.URL.Query().Get("nome"); v != "" {
			filters["nome"] = v
		}
		if v := r.URL.Query().Get("fornecedor_id"); v != "" {
			id, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				filters["fornecedor_id"] = id
			}
		}
		if v := r.URL.Query().Get("codigo_fornecedor"); v != "" {
			filters["codigo_fornecedor"] = v
		}
		// Preço mínimo e máximo
		if v := r.URL.Query().Get("preco_min"); v != "" {
			dec := apd.New(0, 0)
			if _, _, err := dec.SetString(v); err == nil {
				filters["preco_min"] = dec
			}
		}
		if v := r.URL.Query().Get("preco_max"); v != "" {
			dec := apd.New(0, 0)
			if _, _, err := dec.SetString(v); err == nil {
				filters["preco_max"] = dec
			}
		}

		// Paginação via middleware
		page := r.Context().Value(middleware.PageKey).(int)
		limit := r.Context().Value(middleware.LimitKey).(int)
		offset := (page - 1) * limit

		// Consulta
		prods, err := pr.Find(filters, limit, offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro na busca: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prods)
	}
}

// DeleteProduto retorna http.HandlerFunc que exclui pelo ID
func DeleteProduto(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := pr.Delete(id); errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao deletar: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// UpdateProduto retorna http.HandlerFunc que atualiza um produto existente
func UpdateProduto(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		// Decodifica body no DTO
		var dto struct {
			Nome              *string `json:"nome"`
			CodigoBarras      *string `json:"codigo_barras"`
			QuantidadeEstoque *int64  `json:"quantidade_estoque"`
			Preco             *string `json:"preco"`
			QtdMinima         *int64  `json:"qtd_minima"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Busca o produto atual para preencher campos não enviados
		produtos, err := pr.Find(nil, 1, int(id)) // ou um método GetByID
		if err != nil || len(produtos) == 0 {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
			return
		}
		p := &produtos[0]

		// Aplica mudanças
		if dto.Nome != nil {
			p.SetNome(*dto.Nome)
		}

		if dto.QuantidadeEstoque != nil {
			p.SetQuantidadeEstoque(*dto.QuantidadeEstoque)
		}
		if dto.Preco != nil {
			p.SetPreco(*dto.Preco)
		}

		// Persiste update
		if err := pr.Update(p); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao atualizar: %v", err), http.StatusInternalServerError)
			return
		}

		// Retorna o objeto atualizado
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}
