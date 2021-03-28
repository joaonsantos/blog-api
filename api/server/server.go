package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// App contains the router vars and database connections
type App struct {
	Handler http.Handler
	DB      *sql.DB
}

// Initialize bootstraps db connections and registers routes
func (a *App) Initialize(c Config) {
	var err error
	a.DB, err = sql.Open("sqlite3", c.DB_DSN) // imported driver through _ import
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	// server.RegisterRoutes(router, router)
	if c.Log {
		a.Handler = handlers.LoggingHandler(os.Stdout, router)
		return
	}

	a.Handler = router
}

func (a *App) Run(addr string) {
	// Good practice to set timeouts, avoids slowloris
	srv := &http.Server{
		Handler:      a.Handler,
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Server started at adress http://%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
