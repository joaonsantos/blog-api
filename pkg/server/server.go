package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/blog-api/pkg/db"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

// getPostInfo fetches a post info by post slug
func getPostInfo(w http.ResponseWriter, r *http.Request, c *pgxpool.Pool) {
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

// getPostsInfo fetches all posts info
func getPostsInfo(w http.ResponseWriter, r *http.Request, c *pgxpool.Pool) {
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

// submitPostInfo submits a post info
func submitPostInfo(w http.ResponseWriter, r *http.Request, c *pgxpool.Pool) {
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

// getPostContent fetches a post content by post slug
func getPostContent(w http.ResponseWriter, r *http.Request, c *pgxpool.Pool) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	dat, err := ioutil.ReadFile("/opt/blog-api/posts/" + slug + "/index.md")
	if err != nil {
		log.Println(err.Error())
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<h1>` + err.Error() + `</h1>`))
	}

	w.Header().Set("Content-Type", "text/markdown")
	w.Write(dat)
}

// postPostContent creates a file with post content given a post slug
func postPostContent(w http.ResponseWriter, r *http.Request, c *pgxpool.Pool) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}

	err = os.MkdirAll("/opt/blog-api/posts/"+slug, 0744)
	if err != nil {
		log.Println(err.Error())
	}

	err = ioutil.WriteFile("/opt/blog-api/posts/"+slug+"/index.md", body, 0644)
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
func RegisterRoutes(r *mux.Router, c *pgxpool.Pool) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<h1>` + `Server is up` + `</h1>`))
	})

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		getPostsInfo(w, r, c)
	}).Methods("GET")
	api.HandleFunc("/post/{slug}", func(w http.ResponseWriter, r *http.Request) {
		getPostInfo(w, r, c)
	}).Methods("GET")
	api.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		submitPostInfo(w, r, c)
	}).Methods("POST")

	api.HandleFunc("/content/{slug}", func(w http.ResponseWriter, r *http.Request) {
		getPostContent(w, r, c)
	}).Methods("GET")
	api.HandleFunc("/content/{slug}", func(w http.ResponseWriter, r *http.Request) {
		postPostContent(w, r, c)
	}).Methods("POST")
	api.HandleFunc("*", notFound)
}
