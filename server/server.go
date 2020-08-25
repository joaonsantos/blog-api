package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blog-api/db"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// getPost fetches a post by post slug
func getPost(w http.ResponseWriter, r *http.Request, c *pgx.Conn) {
  vars := mux.Vars(r)
  slug := vars["slug"]
	post, err := db.GetPost(c, slug)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<h1>` + err.Error() + `</h1>`))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(post)
}

// getPosts fetches all posts
func getPosts(w http.ResponseWriter, r *http.Request, c *pgx.Conn) {
	posts, err := db.GetPosts(c)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<h1>` + err.Error() + `</h1>`))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(posts)
}

// submitPost submits a post
func submitPost(w http.ResponseWriter, r *http.Request, c *pgx.Conn) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}
	var p db.Post
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Println(err.Error())
	}

	err = db.SubmitPost(c, &p)
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<h1>` + err.Error() + `</h1>`))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// notFound handles not found routes
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

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<h1>` + `Server is up` + `</h1>`))
	})

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		getPosts(w, r, c)
	}).Methods("GET")
	api.HandleFunc("/post/{slug}", func(w http.ResponseWriter, r *http.Request) {
		getPost(w, r, c)
	}).Methods("GET")
	api.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		submitPost(w, r, c)
	}).Methods("POST")
	api.HandleFunc("*", notFound)
}
