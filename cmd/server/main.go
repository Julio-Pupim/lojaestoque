package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/julio-pupim/lojaestoque/internal/database"
	handler "github.com/julio-pupim/lojaestoque/internal/handler"
	myMiddleware "github.com/julio-pupim/lojaestoque/internal/middleware"
	"github.com/julio-pupim/lojaestoque/internal/repository"
)

func main() {
	db := database.InitDB()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/clientes", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handler.CriarCliente(db, w, r)
		})
		r.With(myMiddleware.Pagination).Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.BuscarClientes(db, w, r)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handler.DeleteCliente(db, w, r)
		})
		r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
			handler.UpdateCliente(db, w, r)
		})
	})
	r.Route("/produtos", func(r chi.Router) {
		pr := repository.NewProdutoRepository(db)
		r.Post("/", handler.CreateOrAddProduto(pr))
		r.With(myMiddleware.Pagination).Get("/", handler.SearchProdutos(pr))
		r.Delete("/{id}", handler.DeleteProduto(pr))
		r.Patch("/{id}", handler.UpdateProduto(pr))
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
