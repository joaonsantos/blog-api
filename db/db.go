package db

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

type post struct {
	PostID int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// GetPosts queries the database for posts and returns them as json
func GetPosts(c *pgx.Conn) ([]byte, error) {
	p := []post{}

	rows, err := c.Query(context.Background(), "select * from posts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var title, body string

		err := rows.Scan(&id, &title, &body)
		if err != nil {
			return nil, err
		}

		p = append(p, post{PostID: id, Title: title, Body: body})
	}

	data, err := json.Marshal(p)

	return data, err
}
