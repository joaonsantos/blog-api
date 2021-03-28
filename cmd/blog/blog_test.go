package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blog-api/api/server"
)

var a server.App

const tableCreationStmt = `create table if not exists posts (
	id               text         not null,
	title            text         not null,
	body             text         not null,
	summary          text         not null,
	author           text         not null,
	readTime         integer      not null,
	dateModified     integer      not null,
	constraint posts_pkey primary key (id)
  );`

func initTables() {
	if _, err := a.DB.Exec(tableCreationStmt); err != nil {
		log.Fatal(err)
	}
}

func clearTables() {
	a.DB.Exec("DELETE FROM posts")
}

func TestMain(m *testing.M) {
	a.Initialize(server.Config{DB_DSN: "/tmp/blog.db", Log: false})

	initTables()
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func TestEmptyPosts(t *testing.T) {
	clearTables()

	req := httptest.NewRequest("GET", "/api/v1/posts", nil)
	rr := httptest.NewRecorder()
	a.Handler.ServeHTTP(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code %v. Got %v.", res.StatusCode, http.StatusOK)
	}
}

func TestNonExistentPost(t *testing.T) {
	clearTables()

	req := httptest.NewRequest("GET", "/api/v1/post/test", nil)
	rr := httptest.NewRecorder()
	a.Handler.ServeHTTP(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected response code %v. Got %v.", res.StatusCode, http.StatusOK)
	}

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["error"] != "post does not exist" {
		t.Errorf("Expected the 'error' key of the response to be set to 'post does not exist'.")
	}
}

func TestCreatePost(t *testing.T) {
	clearTables()

	postJSON := []byte(`{
		"title": "Programming is More Than Syntax",
		"summary": "What makes up a programming language.",
		"body": "This is the content of the post."
		"author": "João Santos"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/post/test", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	a.Handler.ServeHTTP(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response code %v. Got %v.", res.StatusCode, http.StatusOK)
	}

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["id"] != "programming-is-more-than-syntax" {
		t.Errorf("Expected the post id to be 'programming-is-more-than-syntax'. Got '%v'.", m["id"])
	}

	if m["readTime"] != 0 {
		t.Errorf("Expected the post read time to be '0'. Got '%v'.", m["readTime"])
	}
}
