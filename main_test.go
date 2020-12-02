package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
	"github.com/pkg/errors"
)

var idOne = int64(1)
var idTwo = int64(2)
var repositoryNameOne = "Repository One"
var repositoryNameTwo = "Repository Two"
var projectNameOne = "Project One"
var projectNameTwo = "Project Two"

// Mocks
var mockBoards = []*trello.Board{
	{
		ID:   "1",
		Name: "Test One",
	},
	{
		ID:   "2",
		Name: "Test Two",
	},
}

var mockRepositories = []*github.Repository{
	{
		ID:   &idOne,
		Name: &repositoryNameOne,
	},
	{
		ID:   &idTwo,
		Name: &repositoryNameTwo,
	},
}

var mockProjects = []*github.Project{
	{
		ID:   &idOne,
		Name: &projectNameOne,
	},
	{
		ID:   &idTwo,
		Name: &projectNameTwo,
	},
}

type mockTrelloClient struct {
}

func (m *mockTrelloClient) GetMember(memberID string, args trello.Arguments) (member *trello.Member, err error) {
	return &trello.Member{
		Username: "testUser",
	}, nil
}

func (m *mockTrelloClient) GetMyBoards(args trello.Arguments) (boards []*trello.Board, err error) {
	return mockBoards, nil
}

type mockGithubAuthenticator struct {
	out io.Writer
}

func (g *mockGithubAuthenticator) AuthenticateGithub(scopes []github.Scope) (token string, err error) {
	fmt.Fprintln(g.out, "Requesting device and user verification codes from GitHub...")
	fmt.Fprintf(g.out, "Please go to https://test.com and enter the code TEST.\n")
	return "", nil
}

type mockUsersService struct {
}

var username = "testUser"

func (m *mockUsersService) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return &github.User{
		Login: &username,
	}, nil, nil
}

type mockRepositoriesService struct {
}

func (m *mockRepositoriesService) List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	if user == "" {
		return nil, nil, errors.Errorf("User must be supplied")
	}
	return mockRepositories, nil, nil
}

func (m *mockRepositoriesService) ListProjects(ctx context.Context, owner, repo string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error) {
	return mockProjects, nil, nil
}

func TestConfigure(t *testing.T) {
	selections := []string{"token", "0", "0", "0", "n", "y"}
	in := bytes.NewBufferString(strings.Join(selections, "\n"))

	out := bytes.Buffer{}

	cliFactory := clientFactory{
		newTrelloClient: func(token string) trelloClient {
			return &mockTrelloClient{}
		},
		newGithubClient: func(token string) githubClient {
			return githubClient{
				users:        &mockUsersService{},
				repositories: &mockRepositoriesService{},
			}
		},
		newGithubAuthenticator: func() githubAuthenticator {
			return &mockGithubAuthenticator{
				out: &out,
			}
		},
	}

	// Return Configuration struct
	configure(in, &out, cliFactory)

	/*
	   Lines like:
	   Board: Selected 'Test One'.

	   Should be split, but have to be one line in test because of the way Fscan
	   seems to treat stdin and buffered string differently.

	   The program works properly when run.
	*/
	wantOut := fmt.Sprintf(`# Configuration
No configuration file found in this directory. Follow the prompts to configure the tool.

## Trello authentication
In your browser, open https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%v
When the authentication process is complete, enter the token you receied here.
Token: Successfully authenticated with Trello as testUser.

## GitHub authentication
Requesting device and user verification codes from GitHub...
Please go to https://test.com and enter the code TEST.
Successfully authenticated with GitHub as testUser.

## Card mapping
### Trello board
Fetching Trello boards...
Select the Trello board you want to import from:
[0] Test One (ID: 1)
[1] Test Two (ID: 2)
Board: Selected 'Test One'.

### GitHub repository
Fetching GitHub repositories...
Select the GitHub repository whose projects you want to export to:
[0] Repository One (ID: 1)
[1] Repository Two (ID: 2)
Repository: Selected 'Repository One'.

### GitHub project
Fetching GitHub projects...
Select the GitHub project you want to export to:
[0] Project One (ID: 1)
[1] Project Two (ID: 2)
Project: Selected 'Project One'.

## Confirm
Your configuration is:
Trello account: testUser
Trello board: Test One (ID: 1)
GitHub account: testUser
GitHub repository: Repository One (ID: 1)
GitHub project: Project One (ID: 1)

Would you like to save this configuration to a 'config.yaml' file for next time?
(y/n) 

Would you like to proceed with the import using these settings?
(y/n) 

Configuration complete.
`, apikey)
	assertParagraphMatch(out.String(), wantOut, t)
}

// Utis
func assertNotNil(got interface{}, t *testing.T) {
	t.Helper()
	if got == nil {
		t.Errorf("Got unexpected nil: %v", got)
	}
}

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

func assertParagraphMatch(got, want string, t *testing.T) {
	t.Helper()
	gotStrings := strings.Split(got, "\n")
	wantStrings := strings.Split(want, "\n")
	for i, gotString := range gotStrings {
		if gotString != wantStrings[i] {
			t.Errorf("Line mismatch: \nGot: \t%v\nWant: \t%v", gotString, wantStrings[i])
			return
		}
	}
}
