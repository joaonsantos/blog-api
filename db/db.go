package db

import (
	"context"
	"encoding/json"
	"time"
  "strings"
  "net/url"

	"github.com/jackc/pgx/v4"
)

type Post struct {
	PostID string    `json:"id"`
	Title  string    `json:"title"`
	Body   string    `json:"body"`
	Date   time.Time `json:"date"`
}

// GetPosts queries the database for posts and returns them as json
func GetPosts(c *pgx.Conn) ([]byte, error) {
	p := []Post{}

	rows, err := c.Query(context.Background(), "select * from posts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, title, body string
		var date time.Time

		err := rows.Scan(&id, &title, &body, &date)
		if err != nil {
			return nil, err
		}

		p = append(p, Post{PostID: id, Title: title, Body: body, Date: date})
	}

	data, err := json.Marshal(p)

	return data, err
}

func formatID(s string) string {
  fid := strings.ReplaceAll(s, " ", "-")
  fid = url.PathEscape(fid)

  return strings.ToLower(fid)
}

// SubmitPost writes a post to the database
func SubmitPost(c *pgx.Conn, p *Post) error {
	title := p.Title
	body := p.Body
	date := time.Now()
  id := formatID(title)

	_, err := c.Exec(context.Background(), "insert into posts values ($1,$2,$3,$4)", id, title, body, date)
	return err
}
