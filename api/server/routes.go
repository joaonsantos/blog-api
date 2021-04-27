package server

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaonsantos/blog-api/pkg/posts"
)

func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	sendResponse(w, code, map[string]string{"error": message})
}

func sendResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) createPost(w http.ResponseWriter, r *http.Request) {
	var p posts.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.CreatePost(a.DB); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)
}

func (a *App) getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Expected to receive url variable 'id'")
		return
	}

	p := posts.Post{ID: id}
	if err := p.GetPost(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			sendErrorResponse(w, http.StatusNotFound, "Post does not exist")
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (a *App) patchPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Expected to receive url variable 'id'")
		return
	}

	var p posts.Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = id
	if err := p.PatchPost(a.DB); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (a *App) getPosts(w http.ResponseWriter, r *http.Request) {
	const postLimit = 100
	posts, err := posts.GetPosts(a.DB, 0, postLimit)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, posts)
}

func (a *App) RegisterRoutes() {
	a.Router.HandleFunc("/api/v1/post", a.createPost).Methods("POST")
	a.Router.HandleFunc("/api/v1/post/{id:[0-9A-Za-z-]+}", a.getPost).Methods("GET")
	a.Router.HandleFunc("/api/v1/post/{id:[0-9A-Za-z-]+}", a.patchPost).Methods("PATCH")
	a.Router.HandleFunc("/api/v1/posts", a.getPosts).Methods("GET")
}
