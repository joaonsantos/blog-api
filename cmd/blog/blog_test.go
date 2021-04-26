package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joaonsantos/blog-api/api/server"
)

var a server.App

const tableCreationStmt = `create table if not exists posts (
	id           text         not null,
	title        varchar(256) not null,
	body         text         not null,
	summary      varchar(256) not null,
	author       varchar(128) not null,
	readTime     integer      not null,
	createDate   integer      not null,
	constraint   posts_pkey   primary key (id)
  );`

func initTables() {
	if _, err := a.DB.Exec(tableCreationStmt); err != nil {
		log.Fatal(err)
	}
}

func clearTables() {
	a.DB.Exec("delete from posts")
}

func addPosts(count int) {
	for i := 0; i < count; i++ {
		a.DB.Exec(`insert into posts(id, title, body, summary, author, readTime)
		values($1,$2,$3,$4,$5,$6,$7)`,
			fmt.Sprintf("test-%v", i),
			fmt.Sprintf("Test %v", i),
			fmt.Sprintf("Test Content %v", i),
			fmt.Sprintf("Summary %v", i),
			"Test Author",
			1+i,
			i,
		)
	}
}

func doRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseError(t *testing.T, payload map[string]interface{}) {
	if err, ok := payload["error"]; ok {
		t.Errorf("error '%v'.", err)
	}
}

func checkResponseCode(t *testing.T, expected int, actual int) bool {
	isErrorCode := actual != expected
	if isErrorCode {
		t.Errorf("Expected response code %v. Got %v.", expected, actual)
	}

	return isErrorCode
}

func TestMain(m *testing.M) {
	a = server.NewApp(&server.Config{DB_DSN: "/tmp/blog.db", Log: false})

	initTables()
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func TestCreatePost(t *testing.T) {
	clearTables()

	postJSON := []byte(`{
		"title": "Programming is More Than Syntax",
		"summary": "What makes up a programming language.",
		"body": "This is the content of the post.",
		"author": "João Santos"
	}`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/post",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	isErrorCode := checkResponseCode(t, http.StatusCreated, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)

	if isErrorCode {
		checkResponseError(t, m)
	}

	if id := m["id"]; id != "programming-is-more-than-syntax" {
		t.Errorf("expected the post id to be 'programming-is-more-than-syntax'. Got '%v'.", id)
	}

	if readTime := m["readTime"]; readTime != 1.0 {
		t.Errorf("expected the post read time to be '1'. Got '%v'.", readTime)
	}

	if _, ok := m["createDate"]; !ok {
		t.Errorf("expected response to contain a 'createDate' field.")
	}
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

func TestGetPosts(t *testing.T) {
	clearTables()
	addPosts(2)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/posts/",
		nil,
	)

	rr := doRequest(req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var l []map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &l)

	if s := len(l); s != 2 {
		t.Errorf("Expected to get two posts. Got '%v'.", s)
	}

	for i := range l {
		m := l[i]

		if m["id"] != fmt.Sprintf("test-%v", i) {
			t.Errorf("Expected the post id to be 'test-1'. Got '%v'.", m["id"])
		}

		if m["readTime"] != 1 {
			t.Errorf("Expected the post read time to be '1'. Got '%v'.", m["readTime"])
		}
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

func TestGetPost(t *testing.T) {
	clearTables()
	addPosts(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/post/test-1",
		nil,
	)

	rr := doRequest(req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["id"] != "test-1" {
		t.Errorf("Expected the post id to be 'test-1'. Got '%v'.", m["id"])
	}

	if m["readTime"] != 1 {
		t.Errorf("Expected the post read time to be '1'. Got '%v'.", m["readTime"])
	}
}

func TestUpdatePost(t *testing.T) {
	clearTables()
	addPosts(1)

	postJSON := []byte(`{
		"summary": "What makes up a test.",
		"body": "This is the content of the test post."
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/post/test-1",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	checkResponseCode(t, http.StatusOK, rr.Code)

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["id"] != "test-1" {
		t.Errorf("Expected the post id to be 'test-1'. Got '%v'.", m["id"])
	}

	if m["summary"] != "What makes up a test." {
		t.Errorf("Expected the post id to be 'What makes up a test.'. Got '%v'.", m["summary"])
	}

	if m["body"] != "This is the content of the test post." {
		t.Errorf("Expected the post id to be 'This is the content of the test post.'. Got '%v'.", m["body"])
	}
}
