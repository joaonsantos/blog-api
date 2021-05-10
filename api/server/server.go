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

// Config contains the configs for the app
type Config struct {
	DB_DSN string
	Log    bool
}

// App contains the router vars and database connections
type App struct {
	Router *mux.Router
	DB     *sql.DB
	Config *Config
}

// NewApp bootstraps db connections and registers routes
func NewApp(c *Config) App {
	a := App{Config: c}

	var err error
	a.DB, err = sql.Open("sqlite3", c.DB_DSN) // imported driver through _ import
	if err != nil {
		log.Fatalf("error initializing db: %v", err)
	}

	err = a.DB.Ping()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}

	a.Router = mux.NewRouter()
	a.RegisterRoutes()

	return a
}

// Run starts a server running on addr, eg. addr=':8080'
func (a *App) Run(addr string) {
	h := handlers.LoggingHandler(os.Stdout, a.Router)

	// if log option is set to false, simply use mux router
	if !a.Config.Log {
		h = a.Router
	}

	// Good practice to set timeouts, avoids slowloris
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server started at adress http://%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
