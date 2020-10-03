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
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Author    string    `json:"author"`
	ReadTime  int       `json:"readTime"`
	Date      int64     `json:"dateModified"`
}

// GetPost queries the database for a post info and returns it as json
func GetPost(c *pgx.Conn, s string) ([]byte, error) {
	p := []Post{}

	rows, err := c.Query(context.Background(), "select slug, title, summary, author, readTime, dateModified from posts where slug=$1", s)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var slug, title, summary, author string
    var readTime int
		var date int64

		err := rows.Scan(&slug, &title, &summary, &author, &readTime, &date)
		if err != nil {
			return nil, err
		}

    p = append(p, Post{Slug: slug, Title: title, Summary: summary, Author: author, ReadTime: readTime, Date: date})
	}

	data, err := json.Marshal(p)

	return data, err
}

// GetPosts queries the database for posts info and returns them as json
func GetPosts(c *pgx.Conn) ([]byte, error) {
	p := []Post{}

	rows, err := c.Query(context.Background(), "select slug, title, summary, author, readTime, dateModified from posts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var slug, title, summary, author string
    var readTime int
		var date int64

		err := rows.Scan(&slug, &title, &summary, &author, &readTime, &date)
		if err != nil {
			return nil, err
		}

    p = append(p, Post{Slug: slug, Title: title, Summary: summary, Author: author, ReadTime: readTime, Date: date})
	}

	data, err := json.Marshal(p)

	return data, err
}

func genSlug(s string) string {
  slug := strings.ReplaceAll(s, " ", "-")
  slug = url.PathEscape(slug)

  return strings.ToLower(slug)
}

// SubmitPost writes a post info to the database
func SubmitPost(c *pgx.Conn, p *Post) error {
	title := p.Title
	summary := p.Summary
  author := p.Author

	readTime := p.ReadTime
  date := time.Now().Unix()
  slug := genSlug(title)

	_, err := c.Exec(context.Background(), "insert into posts values ($1,$2,$3,$4,$5,$6)", slug, title, summary, author, readTime, date)
	return err
}
