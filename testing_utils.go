package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/v32/github"
)

func assertNoError(err error, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}
}

func assertStringMatch(got string, want string, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %s, want %s", got, want)
	}
}

func assertStructsMatch(got interface{}, want interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got %v, want %v", got, want)
	}
}

const (
	baseURLPath = "/api-v3"
)

func mockGithubClient() (client *github.Client, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that. See issue #752.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client = github.NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url
	client.UploadURL = url

	return client, mux, server.Close
}
