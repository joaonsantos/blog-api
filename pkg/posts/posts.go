package posts

import (
	"database/sql"
	"strings"
	"time"

	"github.com/joaonsantos/blog-api/pkg/math"
)

type Posts []Post

// TODO split into PostInfo and PostContent structs

type Post struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Summary    string `json:"summary"`
	Author     string `json:"author"`
	ReadTime   int    `json:"readTime"`
	CreateDate int64  `json:"createDate"`
}

func (p *Post) prepareNewPost() {
	titleWords := strings.Split(strings.ToLower(p.Title), " ")
	p.ID = strings.Join(titleWords, "-")
	p.CreateDate = time.Now().Unix()
}

func (p *Post) CreatePost(db *sql.DB) error {
	p.prepareNewPost()

	row := db.QueryRow(
		`insert into posts(id, title, body, summary, author, readTime, createDate)
		values($1, $2, $3, $4, $5, $6, $7) returning id, readTime, createDate`,
		p.ID,
		p.Title,
		p.Body,
		p.Summary,
		p.Author,
		p.ReadTime,
		p.CreateDate,
	)

	return row.Scan(&p.ID, &p.ReadTime, &p.CreateDate)
}

func (p *Post) SubmitPostContent(db *sql.DB) error {
	row := db.QueryRow(
		`update posts
		set body=$1, readTime=$2
		where id=$3
		returning id, readTime`,
		p.Body,
		p.ReadTime,
		p.ID,
	)

	return row.Scan(&p.ID, &p.ReadTime)
}

func (p *Post) GetPost(db *sql.DB) error {
	row := db.QueryRow(
		`select title, body, summary, author, readTime, createDate from posts where id=$1`,
		p.ID,
	)

	return row.Scan(
		&p.Title,
		&p.Body,
		&p.Summary,
		&p.Author,
		&p.ReadTime,
		&p.CreateDate,
	)
}

func (p *Post) GetPostContent(db *sql.DB) error {
	row := db.QueryRow(
		`select body from posts where id=$1`,
		p.ID,
	)

	return row.Scan(
		&p.Body,
	)
}

func (p *Post) PatchPost(db *sql.DB) error {
	p.ReadTime = CalculatePostReadTime(p.Body)

	_, err := db.Exec(
		`update posts set body=$1, summary=$2, readTime=$3 where id=$4`,
		p.Body,
		p.Summary,
		p.ReadTime,
		p.ID,
	)

	return err
}

func GetPosts(db *sql.DB, start, count int) (Posts, error) {
	rows, err := db.Query(
		`select id, title, body, summary, author, readTime, createDate from posts limit $1 offset $2`,
		count,
		start,
	)
	if err != nil {
		return nil, err
	}

	posts := Posts{}
	for rows.Next() {
		var p Post
		rows.Scan(
			&p.ID,
			&p.Title,
			&p.Body,
			&p.Summary,
			&p.Author,
			&p.ReadTime,
			&p.CreateDate,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func CalculatePostReadTime(body string) int {
	// TODO
	// quick workaround, need to parse markdown and take into account
	// code, images, etc.
	modifier := 1.0
	if strings.Contains(body, "```") {
		modifier = 0.4
	}
	return int(float64(math.Max(1, int(len(body)/200))) * modifier)
}
