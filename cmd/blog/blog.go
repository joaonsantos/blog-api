package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blog-api/pkg/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("PG_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	router := initialize(dbpool)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    "0.0.0.0:8000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server started at adress http://%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func initialize(dbpool *pgxpool.Pool) *mux.Router {
	// Create a new router and register routes
	router := mux.NewRouter()
	server.RegisterRoutes(router, dbpool)

	return router
}
