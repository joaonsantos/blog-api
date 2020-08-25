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
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Summary string    `json:"summary"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
}

// GetPost queries the database for a post and returns it as json
func GetPost(c *pgx.Conn, s string) ([]byte, error) {
	p := []Post{}

	rows, err := c.Query(context.Background(), "select * from posts where slug=$1", s)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var slug, title, body, summary, author string
		var date time.Time

		err := rows.Scan(&slug, &title, &body, &summary, &author, &date)
		if err != nil {
			return nil, err
		}

    p = append(p, Post{Slug: slug, Title: title, Body: body, Summary: summary, Author: author, Date: date})
	}

	data, err := json.Marshal(p)

	return data, err
}

// GetPosts queries the database for posts and returns them as json
func GetPosts(c *pgx.Conn) ([]byte, error) {
	p := []Post{}

	rows, err := c.Query(context.Background(), "select slug, title, summary, author, date from posts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var slug, title, summary, author string
		var date time.Time

		err := rows.Scan(&slug, &title, &summary, &author, &date)
		if err != nil {
			return nil, err
		}

    p = append(p, Post{Slug: slug, Title: title, Summary: summary, Author: author, Date: date})
	}

	data, err := json.Marshal(p)

	return data, err
}

func genSlug(s string) string {
  slug := strings.ReplaceAll(s, " ", "-")
  slug = url.PathEscape(slug)

  return strings.ToLower(slug)
}

// SubmitPost writes a post to the database
func SubmitPost(c *pgx.Conn, p *Post) error {
	title := p.Title
	body := p.Body
	summary := p.Summary
  author := p.Author

	date := time.Now()
  slug := genSlug(title)

	_, err := c.Exec(context.Background(), "insert into posts values ($1,$2,$3,$4,$5,$6)", slug, title, body, summary, author, date)
	return err
}
