package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Handle("/*", http.FileServer(http.Dir("./frontend")))

	r.Route("/clientes", func(r chi.Router) {
		cr := repository.NewClienteRepository(db)
		r.Post("/", handler.CriarCliente(cr))
		r.With(myMiddleware.Pagination).Get("/", handler.BuscarClientes(cr))
		r.Delete("/{id}", handler.DeleteCliente(cr))
		r.Patch("/{id}", handler.UpdateCliente(cr))
		r.Get("/{id}", handler.GetClienteById(cr))
	})
	r.Route("/produtos", func(r chi.Router) {
		pr := repository.NewProdutoRepository(db)
		r.Post("/", handler.CreateOrAddProduto(pr))
		r.With(myMiddleware.Pagination).Get("/", handler.SearchProdutos(pr))
		r.Delete("/{id}", handler.DeleteProduto(pr))
		r.Patch("/{id}", handler.UpdateProduto(pr))
	})
	r.Route("/vendas", func(r chi.Router) {
		vr := repository.NewVendasRepository(db)
		r.Post("/", handler.CriarVenda(vr))
		r.With(myMiddleware.Pagination).Get("/", handler.BuscarVenda(vr))
		r.Delete("/{id}", handler.DeletarVenda(vr))
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
