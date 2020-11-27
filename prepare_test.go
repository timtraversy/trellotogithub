package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	wantOut := fmt.Sprintf(`Select the GitHub project you want to export to:
(0) %v
(1) %v
`, mockProjects[0].Name, mockProjects[1].Name)
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
