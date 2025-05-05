package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/julio-pupim/lojaestoque/internal/domain"
	"github.com/julio-pupim/lojaestoque/internal/middleware"
	"github.com/julio-pupim/lojaestoque/internal/repository"
)

func CriarVenda(vr *repository.VendasRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var venda domain.Sale
		if err := json.NewDecoder(r.Body).Decode(&venda); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Erro ao decodificar JSON")
			return
		}

		err := vr.SalvarVenda(&venda)
		if err != nil {
			log.Printf("Erro ao executar query: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Erro ao inserir cliente")
			return
		}
		RespondCreated(w, venda)
	}
}

func BuscarVenda(vr *repository.VendasRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := make(map[string]any)
		if v := r.URL.Query().Get("nome_cliente"); v != "" {
			filters["nome"] = v
		}
		if v := r.URL.Query().Get("nome_produto"); v != "" {
			filters["nome_produto"] = v
		}
		if v := r.URL.Query().Get("total"); v != "" {
			filters["total"] = v
		}
		if v := r.URL.Query().Get("data_venda"); v != "" {
			filters["data_venda"] = v
		}
		if v := r.URL.Query().Get("data_pagamento"); v != "" {
			filters["data_pagamento"] = v
		}
		if v := r.URL.Query().Get("status_pagamento"); v != "" {
			filters["status_pagamento"] = v
		}
		page := r.Context().Value(middleware.PageKey).(int)
		limit := r.Context().Value(middleware.LimitKey).(int)
		offset := (page - 1) * limit
		vendas, err := vr.BuscarVendas(filters, limit, offset)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("Erro ao buscar vendas: %v", err))
			return
		}
		RespondOK(w, vendas)
	}
}
func DeletarVenda(vr *repository.VendasRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}
		err = vr.DeletarVenda(id)
		if err != nil {
			log.Printf("Erro ao executar query de delete: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Erro interno")
			return
		}
		RespondNoContent(w)
	}
}
