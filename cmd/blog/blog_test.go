package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joaonsantos/blog-api/api/server"
)

var a server.App

const tableCreationStmt = `
create table if not exists posts (
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
		a.DB.Exec(`
			insert into posts(id, title, body, summary, author, readTime, createDate)
			values($1,$2,$3,$4,$5,$6,$7)`,
			fmt.Sprintf("test-%v", i),
			fmt.Sprintf("Test %v", i),
			fmt.Sprintf("Test Content %v", i),
			fmt.Sprintf("Summary %v", i),
			"Test Author",
			1+i,
			i,
			0,
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
	a = server.NewApp(&server.Config{DB_DSN: "file::memory?mode=memory&cache=shared", Log: false})

	initTables()
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func TestGetNonExistentPosts(t *testing.T) {
	clearTables()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/posts/info",
		nil,
	)

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()

	checkResponseCode(t, http.StatusOK, result.StatusCode)

	resBodyBytes, _ := io.ReadAll(result.Body)
	resBody := string(resBodyBytes)
	if resBody != "[]" {
		t.Errorf("expected to get empty payload, got '%v'", resBody)
	}

	var l []map[string]interface{}
	json.Unmarshal(resBodyBytes, &l)

	if s := len(l); s != 0 {
		t.Errorf("Expected to get 0 posts. Got '%v'.", s)
	}

}

func TestGetPosts(t *testing.T) {
	clearTables()
	addPosts(2)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/posts/info",
		nil,
	)

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()

	checkResponseCode(t, http.StatusOK, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var l []map[string]interface{}
	json.Unmarshal(resBody, &l)

	if s := len(l); s != 2 {
		t.Errorf("Expected to get two posts. Got '%v'.", s)
	}

	for i := range l {
		m := l[i]

		offset := float64(i)
		expectedId := fmt.Sprintf("test-%v", offset)
		if id := m["id"]; id != expectedId {
			t.Errorf("Expected the post id to be 'test-%v'. Got '%v'.", offset, m["id"])
		}
		expectedReadTime := 1.0 + offset
		if readTime := m["readTime"]; readTime != expectedReadTime {
			t.Errorf("Expected the post read time to be '1'. Got '%v'.", m["readTime"])
		}
	}
}

func TestGetNonExistentPost(t *testing.T) {
	clearTables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/post/info/test", nil)
	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()

	checkResponseCode(t, http.StatusNotFound, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if err, ok := m["error"]; ok {
		if err != "Post does not exist" {
			t.Errorf("Expected the error to be 'Post does not exist'. Got '%v'.", err)
		}
	}
}

func TestGetPost(t *testing.T) {
	clearTables()
	addPosts(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/post/info/test-0",
		nil,
	)

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()

	checkResponseCode(t, http.StatusOK, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if id := m["id"]; id != "test-0" {
		t.Errorf("Expected the post id to be 'test-0'. Got '%v'.", id)
	}
	if readTime := m["readTime"]; readTime != 1.0 {
		t.Errorf("expected the post read time to be '1'. Got '%v'.", readTime)
	}
}

func TestGetPostContent(t *testing.T) {
	clearTables()
	addPosts(1)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/post/content/test-0",
		nil,
	)

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()

	checkResponseCode(t, http.StatusOK, result.StatusCode)
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
		"/api/v1/post/info",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()
	isErrorCode := checkResponseCode(t, http.StatusCreated, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if isErrorCode {
		checkResponseError(t, m)
	}

	if id := m["id"]; id != "programming-is-more-than-syntax" {
		t.Errorf("expected the post id to be 'programming-is-more-than-syntax'. Got '%v'.", id)
	}
	if readTime := m["readTime"]; readTime != 0.0 {
		t.Errorf("expected the post read time to be '0'. Got '%v'.", readTime)
	}
	if _, ok := m["createDate"]; !ok {
		t.Errorf("expected response to contain a 'createDate' field.")
	}
}

func TestSubmitPostContent(t *testing.T) {
	post, err := os.ReadFile("./post.md")
	if err != nil {
		t.Fatalf("unable to open post contents file")
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/post/content/programming-is-more-than-syntax",
		bytes.NewBuffer(post),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()
	isErrorCode := checkResponseCode(t, http.StatusOK, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if isErrorCode {
		checkResponseError(t, m)
	}
}

func TestPatchPost(t *testing.T) {
	clearTables()
	addPosts(1)

	postJSON := []byte(`{
		"summary": "What makes up a test.",
		"body": "This is the content of the test post."
	}`)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/post/info/test-0",
		bytes.NewBuffer(postJSON),
	)
	req.Header.Set("Content-Type", "application/json")

	rr := doRequest(req)
	result := rr.Result()
	result.Body.Close()
	checkResponseCode(t, http.StatusOK, result.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/post/info/test-0", nil)
	rr = doRequest(req)
	result = rr.Result()
	defer result.Body.Close()

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if id := m["id"]; id != "test-0" {
		t.Errorf("Expected the post id to be 'test-0'. Got '%v'.", id)
	}
	if summary := m["summary"]; summary != "What makes up a test." {
		t.Errorf("Expected the post summary to be 'What makes up a test.'. Got '%v'.", summary)
	}
	if body := m["body"]; body != "This is the content of the test post." {
		t.Errorf("Expected the post body to be 'This is the content of the test post.'. Got '%v'.", body)
	}
}

func TestEmptyPosts(t *testing.T) {
	clearTables()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/info", nil)
	rr := doRequest(req)
	result := rr.Result()
	defer result.Body.Close()
	isErrorCode := checkResponseCode(t, http.StatusOK, result.StatusCode)

	resBody, _ := io.ReadAll(result.Body)
	var m map[string]interface{}
	json.Unmarshal(resBody, &m)

	if isErrorCode {
		checkResponseError(t, m)
	}
	if size := len(m); size != 0 {
		t.Errorf("expected response to be empty, got response with %v items", size)
	}
}
