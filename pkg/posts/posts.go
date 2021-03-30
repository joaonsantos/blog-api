package posts

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Post struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Summary    string `json:"summary"`
	Author     string `json:"author"`
	ReadTime   int    `json:"readTime"`
	CreateDate int64  `json:"createDate"`
}

type Posts []Post

// newPost creates a new Post and returns it
func newPost(title, summary, body, author string) Post {
	p := Post{
		Title:   title,
		Body:    body,
		Summary: summary,
		Author:  author,
	}

	titleWords := strings.Split(strings.ToLower(p.Title), " ")
	p.ID = strings.Join(titleWords, "-")
	p.ReadTime = int(len(p.Body) / 200)
	p.CreateDate = time.Now().Unix()

	return p
}

func getPosts(db *sql.DB, start, count int) (Posts, error) {
	return nil, errors.New("not implemented")
}
