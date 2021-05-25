package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaonsantos/blog-api/pkg/posts"
)

func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	sendResponseJSON(w, code, map[string]string{"error": message})
}

func sendResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func sendResponseMarkdown(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "text/markdown")
	w.WriteHeader(code)
	w.Write([]byte(payload))
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

	sendResponseJSON(w, http.StatusCreated, p)
}

func (a *App) submitPostContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Expected to receive url variable 'id'")
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
	}
	defer r.Body.Close()

	p := posts.Post{
		ID:       id,
		Body:     string(content),
		ReadTime: posts.CalculatePostReadTime(string(content)),
	}

	if err := p.SubmitPostContent(a.DB); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	payload := make(map[string]interface{})
	payload["id"] = p.ID
	payload["readTime"] = p.ReadTime
	sendResponseJSON(w, http.StatusOK, payload)
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

	sendResponseJSON(w, http.StatusOK, p)
}

func (a *App) getPostContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		sendErrorResponse(w, http.StatusBadRequest, "Expected to receive url variable 'id'")
		return
	}

	p := posts.Post{ID: id}
	if err := p.GetPostContent(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			sendErrorResponse(w, http.StatusNotFound, "Post does not exist")
		default:
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponseMarkdown(w, http.StatusOK, p.Body)
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

	sendResponseJSON(w, http.StatusOK, p)
}

func (a *App) getPosts(w http.ResponseWriter, r *http.Request) {
	const postLimit = 100
	posts, err := posts.GetPosts(a.DB, 0, postLimit)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponseJSON(w, http.StatusOK, posts)
}

func (a *App) RegisterRoutes() {
	a.Router.HandleFunc("/api/v1/post/info", a.createPost).Methods("POST")
	a.Router.HandleFunc("/api/v1/post/content/{id:[0-9A-Za-z-]+}", a.getPostContent).Methods("GET")
	a.Router.HandleFunc("/api/v1/post/content/{id:[0-9A-Za-z-]+}", a.submitPostContent).Methods("POST")
	a.Router.HandleFunc("/api/v1/post/info/{id:[0-9A-Za-z-]+}", a.getPost).Methods("GET")
	a.Router.HandleFunc("/api/v1/post/info/{id:[0-9A-Za-z-]+}", a.patchPost).Methods("PATCH")
	a.Router.HandleFunc("/api/v1/posts/info", a.getPosts).Methods("GET")
}
