package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
)

func mockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/members/me/boards":
			jsonResponse, _ := json.Marshal(mockBoards)
			w.Write(jsonResponse)
		}
	}))
}

var mockBoards = []trello.Board{
	{
		ID:   "1",
		Name: "Test",
	},
	{
		ID:   "1",
		Name: "Test Two",
	},
}

func mockClient(baseURL string) *trello.Client {
	client := trello.NewClient(apikey, testtoken)
	client.BaseURL = baseURL
	return client
}

func TestSelectBoard(t *testing.T) {
	server := mockServer()
	defer server.Close()
	client := mockClient(server.URL)

	out := &bytes.Buffer{}
	wantOut := fmt.Sprintf(`Select the Trello board you want to import from:
(0) %v
(1) %v
`, mockBoards[0].Name, mockBoards[1].Name)
	selection := 0
	in := bytes.Buffer{}
	in.Write([]byte(fmt.Sprint(selection)))

	gotBoard := selectBoard(client, &in, out)

	assertStringMatch(gotBoard.ID, mockBoards[selection].ID, t)
	assertStringMatch(out.String(), wantOut, t)
}

var idOne = int64(1)
var nameOne = "Project One"
var idTwo = int64(2)
var nameTwo = "Project Two"

type mockUsersClient struct {
}

func (m *mockUsersClient) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return &github.User{}, nil, nil
}

func (m *mockUsersClient) ListProjects(ctx context.Context, user string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error) {
	return mockProjects, nil, nil
}

var mockProjects = []*github.Project{
	{
		ID:   &idOne,
		Name: &nameOne,
	},
	{
		ID:   &idTwo,
		Name: &nameTwo,
	},
}

func TestSelectProject(t *testing.T) {
	out := &bytes.Buffer{}
	wantOut := fmt.Sprintf(`Sxlect the GitHub project you want to export to:
(0) %v
(1) %v
`, *mockProjects[0].Name, *mockProjects[1].Name)
	selection := 0
	in := bytes.Buffer{}
	in.Write([]byte(fmt.Sprint(selection)))

	muc := mockUsersClient{}

	gotProject := selectProject(&muc, &in, out)
	gotID := fmt.Sprint(*gotProject.ID)

	wantProject := mockProjects[selection]
	wantID := fmt.Sprint(*wantProject.ID)

	assertStringMatch(gotID, wantID, t)
	assertStringMatch(out.String(), wantOut, t)
}

// Utis
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
