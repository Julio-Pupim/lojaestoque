package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	domain "github.com/julio-pupim/lojaestoque/internal/domain"
	"github.com/julio-pupim/lojaestoque/internal/middleware"
)

func CriarCliente(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var cliente domain.Cliente

	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Erro ao decodificar JSON")
		return
	}

	cliente.DataCadastro = time.Now().Format("2006-01-02")
	stmt, err := db.Prepare("INSERT INTO clientes (nome, telefone, data_cadastro) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Erro ao preparar query: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Erro interno")
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(cliente.Nome, cliente.Telefone, cliente.DataCadastro)
	if err != nil {
		log.Printf("Erro ao executar query: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Erro ao inserir cliente")
		return
	}

	cliente.ID, _ = result.LastInsertId()
	RespondCreated(w, cliente)
}

func BuscarClientes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// 1. Extrai parâmetros de busca e paginação do contexto e da URL
	nome := r.URL.Query().Get("nome")
	page := r.Context().Value(middleware.PageKey).(int)
	limit := r.Context().Value(middleware.LimitKey).(int)
	offset := (page - 1) * limit

	// 2. Monta a SQL base e argumentos dinâmicos
	baseSQL := "SELECT * FROM clientes"
	var args []any

	if nome != "" {
		baseSQL += " WHERE nome LIKE ?"
		args = append(args, "%"+nome+"%")
	}
	baseSQL += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// 3. Executa a query com tratamento de erro
	rows, err := db.Query(baseSQL, args...)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar clientes")
		return
	}
	defer rows.Close()

	// 4. Varre resultados e popula slice de domain.Cliente
	var clientes []domain.Cliente
	for rows.Next() {
		var c domain.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Telefone, &c.DataCadastro); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao ler resultados")
			return
		}
		clientes = append(clientes, c)
	}
	if err = rows.Err(); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro após varrer resultados")
		return
	}

	// 5. Retorna JSON
	RespondOK(w, clientes)
}

func DeleteCliente(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	result, err := db.Exec("DELETE FROM clientes WHERE id = ?", id)
	if err != nil {
		log.Printf("Erro ao executar query de delete: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Erro interno")
		return
	}

	// Verifica se algum registro foi afetado
	rows, _ := result.RowsAffected()
	if rows == 0 {
		RespondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}

	RespondNoContent(w)
}

func UpdateCliente(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
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

	var sets []string
	var args []any
	if input.Nome != nil {
		sets = append(sets, "nome = ?")
		args = append(args, *input.Nome)
	}
	if input.Telefone != nil {
		sets = append(sets, "telefone = ?")
		args = append(args, *input.Telefone)
	}
	if len(sets) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Nada para atualizar")
		return
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE clientes SET %s WHERE id = ?", strings.Join(sets, ", "))
	res, err := db.Exec(query, args...)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao atualizar cliente")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		RespondWithError(w, http.StatusNotFound, "Cliente não encontrado")
		return
	}

	var updated domain.Cliente
	err = db.QueryRow("SELECT id, nome, telefone, data_cadastro FROM clientes WHERE id = ?", id).
		Scan(&updated.ID, &updated.Nome, &updated.Telefone, &updated.DataCadastro)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar cliente atualizado")
		return
	}

	RespondOK(w, updated)
}
