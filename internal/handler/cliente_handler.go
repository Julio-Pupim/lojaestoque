package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	domain "github.com/julio-pupim/lojaestoque/internal/domain"
	"github.com/julio-pupim/lojaestoque/internal/middleware"
	"github.com/julio-pupim/lojaestoque/internal/repository"
)

func CriarCliente(cr *repository.ClienteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cliente domain.Cliente

		if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Erro ao decodificar JSON")
			return
		}
		cliente.DataCadastro = time.Now().Format("2006-01-02")

		result, err := cr.SalvarCliente(&cliente)
		if err != nil {
			log.Printf("Erro ao executar query: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Erro ao inserir cliente")
			return
		}
		cliente.ID, _ = result.LastInsertId()
		RespondCreated(w, cliente)
	}
}

func BuscarClientes(cr *repository.ClienteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extrai parâmetros de busca e paginação do contexto e da URL
		filters := make(map[string]any)
		if v := r.URL.Query().Get("nome"); v != "" {
			filters["nome"] = v
		}
		if v := r.URL.Query().Get("telefone"); v != "" {
			filters["telefone"] = v
		}
		page := r.Context().Value(middleware.PageKey).(int)
		limit := r.Context().Value(middleware.LimitKey).(int)
		offset := (page - 1) * limit
		clientes, err := cr.BuscarClientes(filters, limit, offset)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprint("Erro ao buscar clientes: %v", err))
			return
		}
		RespondOK(w, clientes)

	}
}

func DeleteCliente(cr *repository.ClienteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}
		err = cr.DeletarCliente(id)
		if err != nil {
			log.Printf("Erro ao executar query de delete: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Erro interno")
			return
		}
		RespondNoContent(w)
	}
}

func UpdateCliente(cr *repository.ClienteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}
		var input struct {
			Nome     *string `json:"nome"`
			Telefone *string `json:"telefone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			RespondWithError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		result, err := cr.AtualizarCliente(id, domain.Cliente{Nome: *input.Nome, Telefone: *input.Telefone})
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao executar update cliente")
		}
		if row, _ := result.RowsAffected(); row == 0 {
			RespondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		}
		cliente, _ := cr.BuscarClientePorId(id)

		RespondOK(w, cliente)

	}
}

func GetClienteById(cr *repository.ClienteRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)

		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		cliente, err := cr.BuscarClientePorId(id)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao executar busca de cliente por id")
		}
		RespondOK(w, cliente)
	}
}
