package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blog-api/server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func main() {
	router := initialize()
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Server started at adress http://%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func initialize() *mux.Router {
	// Create a new router
	router := mux.NewRouter()

	// Create a new connection to our pg database
	conn, err := pgx.Connect(context.Background(), os.Getenv("PG_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	server.RegisterRoutes(router, conn)

	return router
}