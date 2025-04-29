package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "../../migrations/estoque.db")

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao testar conexão com o banco: %v", err)
	}

	return db
}
