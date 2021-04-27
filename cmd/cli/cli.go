package main

import (
	"flag"
	"fmt"
	"os"

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
	db := flag.String("db", ":8080", "the db dsn, eg. \"blog.db\"")
	flag.Parse()

	a := server.NewApp(&server.Config{
		DB_DSN: *db,
		Log:    true,
	})

	if _, err := a.DB.Exec(tableCreationStmt); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize tables")
	}

}
