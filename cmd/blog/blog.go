package main

import (
	"flag"
	"os"

	"github.com/joaonsantos/blog-api/api/server"
)

func main() {
	serverAddr := flag.String("addr", ":8000", "the addr the server listens on, eg. \":8000\"")
	flag.Parse()

	a := server.NewApp(&server.Config{
		DB_DSN: os.Getenv("DB_DSN"),
		Log:    true,
	})

	a.Run(*serverAddr)
}
