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
		}
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON inválido")
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
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Persiste
		if err := pr.Save(p); err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro ao salvar produto: %v", err))
			return
		}

		// Retorna criado/atualizado
		RespondCreated(w, p)
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
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro na busca: %v", err))
			return
		}

		RespondOK(w, prods)
	}
}

// DeleteProduto retorna http.HandlerFunc que exclui pelo ID
func DeleteProduto(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		if err := pr.Delete(id); errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusNotFound, "Produto não encontrado")
			return
		} else if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro ao deletar: %v", err))
			return
		}

		RespondNoContent(w)
	}
}

// UpdateProduto retorna http.HandlerFunc que atualiza um produto existente
func UpdateProduto(pr *repository.ProdutoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		// Decodifica body no DTO
		var dto struct {
			Nome              *string `json:"nome"`
			CodigoFornecedor  *string `json:"codigo_fornecedor"`
			QuantidadeEstoque *int64  `json:"quantidade_estoque"`
			Preco             *string `json:"preco"`
		}
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		// Busca o produto atual para atualizar
		produtos, err := pr.Find(map[string]any{"id": id}, 1, 0)
		if err != nil || len(produtos) == 0 {
			RespondWithError(w, http.StatusNotFound, "Produto não encontrado")
			return
		}
		p := &produtos[0]

		// Aplica mudanças
		if dto.Nome != nil {
			if err := p.SetNome(*dto.Nome); err != nil {
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if dto.CodigoFornecedor != nil {
			p.CodigoFornecedor = *dto.CodigoFornecedor
		}

		if dto.QuantidadeEstoque != nil {
			if err := p.SetQuantidadeEstoque(*dto.QuantidadeEstoque); err != nil {
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		if dto.Preco != nil {
			if err := p.SetPreco(*dto.Preco); err != nil {
				RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		// Valida e persiste update
		if err := p.ValidateAndUpdate(); err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := pr.Update(p); err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro ao atualizar: %v", err))
			return
		}

		RespondOK(w, p)
	}
}
