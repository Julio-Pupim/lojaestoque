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
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	cliente.DataCadastro = time.Now().Format("2006-01-02")
	stmt, err := db.Prepare("INSERT INTO clientes (nome, telefone, data_cadastro) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Erro ao preparar query: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(cliente.Nome, cliente.Telefone, cliente.DataCadastro)
	if err != nil {
		log.Printf("Erro ao executar query: %v", err)
		http.Error(w, "Erro ao inserir cliente", http.StatusInternalServerError)
		return
	}

	cliente.ID, _ = result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
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
		http.Error(w, "Erro ao buscar clientes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 4. Varre resultados e popula slice de domain.Cliente
	var clientes []domain.Cliente
	for rows.Next() {
		var c domain.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Telefone, &c.DataCadastro); err != nil {
			http.Error(w, "Erro ao ler resultados", http.StatusInternalServerError)
			return
		}
		clientes = append(clientes, c)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, "Erro após varrer resultados", http.StatusInternalServerError)
		return
	}

	// 5. Retorna JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func DeleteCliente(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Erro ao converter id para numérico", http.StatusInternalServerError)
		log.Printf("%v", err)
		return
	}
	_, err = db.Exec("DELETE FROM clientes WHERE id = ?", id)

	if err != nil {
		log.Printf("Erro ao preparar query de delete: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	fmt.Printf("id: %d", id)
	w.WriteHeader(http.StatusNoContent)
}
func UpdateCliente(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Erro ao converter id para numérico", http.StatusInternalServerError)
		log.Printf("%v", err)
	}
	var input struct {
		Nome     *string `json:"nome"`
		Telefone *string `json:"telefone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
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
		http.Error(w, "Nada para atualizar", http.StatusBadRequest)
		return
	}
	args = append(args, id)
	query := fmt.Sprintf("UPDATE clientes SET %s WHERE id = ?", strings.Join(sets, ", "))
	res, err := db.Exec(query, args...)
	if err != nil {
		http.Error(w, "Erro ao atualizar cliente", http.StatusInternalServerError)
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	var updated domain.Cliente
	err = db.QueryRow("SELECT id, nome, telefone, data_cadastro FROM clientes WHERE id = ?", id).
		Scan(&updated.ID, &updated.Nome, &updated.Telefone, &updated.DataCadastro)
	if err != nil {
		http.Error(w, "Erro ao buscar cliente atualizado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(updated)
}
