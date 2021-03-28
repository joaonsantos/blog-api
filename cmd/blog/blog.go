package main

import (
	"flag"
	"os"

	"github.com/blog-api/api/server"
)

var serveraddr string

func flagInit() {
	const (
		defaultAddr = ":8080"
		usage       = "the addr the server listens on, eg. \":8080\""
	)

	flag.StringVar(&serveraddr, "addr", defaultAddr, usage)
	flag.StringVar(&serveraddr, "a", defaultAddr, usage+"(shorthand)")
}

func main() {
	flagInit()
	flag.Parse()

	a := new(server.App)
	a.Initialize(server.Config{
		DB_DSN: os.Getenv("DB_DSN"),
		Log:    true,
	})
	a.Run(serveraddr)
}
