package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

type post struct {
	PostID int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func getDbPosts(c *pgx.Conn) ([]byte, error) {
	p := []post{}

	rows, err := c.Query(context.Background(), "select * from posts;")
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

func getPosts(w http.ResponseWriter, r *http.Request, c *pgx.Conn) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	posts, err := getDbPosts(c)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println(string(posts))
	w.Write(posts)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`<h1>Not Found</h1>`))
}

// RegisterRoutes assigns routes to function handlers
func RegisterRoutes(r *mux.Router, c *pgx.Conn) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		getPosts(w, r, c)
	})
	api.HandleFunc("*", notFound)
}
