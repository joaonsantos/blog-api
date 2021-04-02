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

func doRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if actual != expected {
		t.Errorf("Expected response code %v. Got %v.", expected, actual)
	}
}

func TestMain(m *testing.M) {
	a = server.NewApp(&server.Config{DB_DSN: "/tmp/blog.db", Log: false})

	initTables()
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func TestEmptyPosts(t *testing.T) {
	clearTables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	rr := doRequest(req)
	checkResponseCode(t, http.StatusOK, rr.Code)
}

func TestGetNonExistentPost(t *testing.T) {
	clearTables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/post/test", nil)
	rr := doRequest(req)
	checkResponseCode(t, http.StatusNotFound, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if len(m) != 0 {
		t.Errorf("Expected the response to be an empty json. Got '%v'.", m)
	}
}

func TestUpdateNonExistentPost(t *testing.T) {
	clearTables()

	postJSON := []byte(`{
		"summary": "What makes up a test.",
		"body": "This is the content of the test post."
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/post/programming-is-more-than-syntax",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	checkResponseCode(t, http.StatusNotFound, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if len(m) != 0 {
		t.Errorf("Expected the response to be an empty json. Got '%v'.", m)
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

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/post/programming-is-more-than-syntax",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	checkResponseCode(t, http.StatusCreated, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["id"] != "programming-is-more-than-syntax" {
		t.Errorf("Expected the post id to be 'programming-is-more-than-syntax'. Got '%v'.", m["id"])
	}

	if m["readTime"] != 1 {
		t.Errorf("Expected the post read time to be '0'. Got '%v'.", m["readTime"])
	}
}

func TestUpdatePost(t *testing.T) {
	postJSON := []byte(`{
		"summary": "What makes up a test.",
		"body": "This is the content of the test post."
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/post/programming-is-more-than-syntax",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["id"] != "programming-is-more-than-syntax" {
		t.Errorf("Expected the post id to be 'programming-is-more-than-syntax'. Got '%v'.", m["id"])
	}

	if m["summary"] != "What makes up a test." {
		t.Errorf("Expected the post id to be 'What makes up a test.'. Got '%v'.", m["summary"])
	}

	if m["body"] != "This is the content of the test post." {
		t.Errorf("Expected the post id to be 'This is the content of the test post.'. Got '%v'.", m["body"])
	}
}
