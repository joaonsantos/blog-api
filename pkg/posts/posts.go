package posts

import (
	"database/sql"
	"errors"
	"strings"
)

type Post struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Summary  string `json:"summary"`
	Author   string `json:"author"`
	ReadTime int    `json:"readTime"`
	Date     int64  `json:"dateModified"`
}

// FillMeta should be called after creating a post.
//
// It sets the post ID and the ReadTime.
func (p *Post) FillMeta() {
	titleWords := strings.Split(strings.ToLower(p.Title), " ")
	p.ID = strings.Join(titleWords, "-")
	p.ReadTime = int(len(p.Body) / 200)
}

type Posts []Post

func getPosts(db *sql.DB, start, count int) (Posts, error) {
	return nil, errors.New("not implemented")
}
