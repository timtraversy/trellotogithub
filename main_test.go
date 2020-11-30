package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func selectBoard(client *trello.Client, in io.Reader, out io.Writer) *trello.Board {
	boards, _ := client.GetMyBoards(trello.Defaults())
	fmt.Fprintln(out, "Select the Trello board you want to import from:")
	for i, board := range boards {
		fmt.Fprintf(out, "(%v) %v\n", i, board.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d", selection)
	return boards[selection]
}

type repositoriesClient interface {
	List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
}

type usersClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
	ListProjects(ctx context.Context, user string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error)
}

func selectRepository(client repositoriesClient, in io.Reader, out io.Writer) *github.Repository {
	repos, _, _ := client.List(context.Background(), "", nil)
	fmt.Fprintln(out, "Select the GitHub repository whose projects you want to export to:")
	for i, repo := range repos {
		fmt.Fprintf(out, "(%v) %v\n", i, *repo.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d", &selection)
	return repos[selection]
}

func selectProject(client usersClient, in io.Reader, out io.Writer) *github.Project {
	user, _, _ := client.Get(context.Background(), "")
	projects, _, _ := client.ListProjects(context.Background(), user.GetName(), nil)
	fmt.Fprintln(out, "Select the GitHub project you want to export to:")
	for i, project := range projects {
		fmt.Fprintf(out, "(%v) %v\n", i, *project.Name)
	}

	var selection int
	fmt.Fscanf(in, "%d\n", &selection)
	return projects[selection]
}
