package posts

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Posts []Post

type Post struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Summary    string `json:"summary"`
	Author     string `json:"author"`
	ReadTime   int    `json:"readTime"`
	CreateDate int64  `json:"createDate"`
}

func (p *Post) GetPost(db *sql.DB) error {
	return db.QueryRow("select * from posts where id=$1", p.ID).Scan(
		&p.ID,
		&p.Title,
		&p.Body,
		&p.Summary,
		&p.Author,
		&p.ReadTime,
		&p.CreateDate,
	)
}

func (p *Post) UpdatePost(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE products SET body=$1, summary=$2 , readTime=$3 WHERE id=$4",
		p.Body,
		p.Summary,
		p.ReadTime,
		p.ID,
	)

	return err
}

func (p *Post) CreatePost(db *sql.DB) error {
	return errors.New("not implemented")
}

// newPost creates a new Post and returns it
// TODO remove
func NewPost(title, summary, body, author string) Post {
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

func GetPosts(db *sql.DB, start, count int) (Posts, error) {
	return nil, errors.New("not implemented")
}
