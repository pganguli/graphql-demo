package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pganguli/hnews/graph"
	"github.com/pganguli/hnews/internal/auth"
	db "github.com/pganguli/hnews/internal/db"
	"github.com/pganguli/hnews/internal/links"
	"github.com/pganguli/hnews/internal/users"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

func panicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("Test panic")
	}
}

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db.Open()
	err := db.DB.AutoMigrate(&links.Link{}, &users.User{})
	if err != nil {
		log.Fatal(err)
	}

	gql_handler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.Authenticate)

	router.Use(middleware.Heartbeat("/ping"))
	router.Handle("/panic", panicHandler())
	router.Handle("/graphql", gql_handler)
	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":"+port, router))
}
