package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joaonsantos/blog-api/api/server"
)

const tableCreationStmt = `create table if not exists posts (
	id           text         not null,
	title        varchar(256) not null,
	body         text         not null,
	summary      varchar(256) not null,
	author       varchar(128) not null,
	readTime     integer      not null,
	createDate   integer      not null,
	constraint   posts_pkey   primary key (id)
  );`

func main() {
	var dsn string
	flag.StringVar(&dsn, "db", "", "the db dsn, eg. \"file:blog.db?cache=shared\"")
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

	if _, err := a.DB.Exec(tableCreationStmt); err != nil {
		fmt.Fprintf(os.Stderr, "error starting server: %v\n", err)
		os.Exit(1)
	}

	dsnPartial := strings.Split(dsn, ":")[1]

	dbName := dsnPartial
	if strings.Contains(dsnPartial, "?") {
		dbName = strings.Split(dsnPartial, "?")[0]
	}

	fmt.Printf("✨ created db '%v' successfully \n\n", dbName)
}
