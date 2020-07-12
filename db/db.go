package db

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

type Post struct {
	PostID int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// GetPosts queries the database for posts and returns them as json
func GetPosts(c *pgx.Conn) ([]byte, error) {
	p := []Post{}

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

		p = append(p, Post{PostID: id, Title: title, Body: body})
	}

	data, err := json.Marshal(p)

	return data, err
}


// SubmitPost writes a post to the database
func SubmitPost(c *pgx.Conn, p *Post) error {
  title := p.Title
  body := p.Body

	_, err := c.Exec(context.Background(), "insert into posts(title,body) values ($1,$2)", title, body)
  return err
}
