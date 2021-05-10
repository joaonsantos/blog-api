package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joaonsantos/blog-api/api/server"
)

func main() {
	var dsn string
	flag.StringVar(&dsn, "db", "", "the db dsn, eg. \"file:blog.db?cache=shared\"")
	serverAddr := flag.String("addr", ":8080", "the addr the server listens on, eg. \":8080\"")
	flag.Parse()

	if dsn == "" {
		err := "missing required -db flag, run with -help flag to check usage"
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}

	a := server.NewApp(&server.Config{
		DB_DSN: dsn,
		Log:    true,
	})

	a.Run(*serverAddr)
}
