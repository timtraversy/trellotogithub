package main

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
)

var idOne = int64(1)
var nameOne = "Project One"
var idTwo = int64(2)
var nameTwo = "Project Two"

// Mocks
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

type mockTrelloClient struct {
}

type mockGithubClient struct {
}

var clientFactory = ClientFactory{
	trelloClientFactory: func(token string) TrelloClient {
		return mockTrelloClient{}
	},
	githubClientFactory: func(token string) GithubClient {
		return mockGithubClient{}
	},
}

func TestConfigure(t *testing.T) {
	in := bytes.Buffer{}
	out := bytes.Buffer{}

	in.WriteString("0")
	in.WriteString("0")
	in.WriteString("n")

	// Return Configuration struct
	configure(&in, &out, clientFactory)

	wantOut := `# Configuration
No configuration file found in this directory. Follow the prompts to configure the tool.

## Trello authentication:
In your browser, open https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%v\n
When the authentication process is complete, enter the token you receied here.
Token:
Successfully authenticated with Trello as testUser.

## GitHub authentication:
Requesting device and user verification codes from GitHub...
Please go to https://test.com and enter the code TEST.
Successfully authenticated with GitHub.

## Card mapping:
### Trello board
Fetching Trello boards...
Select the Trello board you want to import from:
[0] Test One (ID: 1)
[1] Test Two (ID: 2)
Board:
Selected 'Test One'.

### GitHub repository 
Fetching GitHub repositories:
Select the GitHub repository whose projects you want to export to:
[0] Repository One (ID: 1)
[1] Repository Two (ID: 2)
Selected 'Repository One'.

### GitHub project 
Fetching GitHub projects:
Select the GitHub project you want to export to:
[0] Project One (ID: 1)
[1] Project Two (ID: 2)
Selected 'Project One'.

## Confirm
Your configuration is:
Trello account: testUser
Trello board: Test One (ID: 1)
GitHub account: testUser
GitHub repository: Repository One (ID: 1)
GitHub project: Project One (ID: 1)

Would you like to proceed with the import using these settings?
(y/n)

Would you like to save this configuration to a 'config.yaml' file for next time?
(y/n) 

Configuration complete.
`
	assertParagraphMatch(out.String(), wantOut, t)
}

type mockUsersClient struct {
}

func (m *mockUsersClient) Get(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return &github.User{}, nil, nil
}

func (m *mockUsersClient) ListProjects(ctx context.Context, user string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error) {
	return mockProjects, nil, nil
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
	fmt.Println(want)
	t.Helper()
	gotStrings := strings.Split(got, "\n")
	wantStrings := strings.Split(want, "\n")

	if len(gotStrings) != len(wantStrings) {
		t.Errorf("Want paragraph with %v lines, got paragraph with %v lines", len(wantStrings), len(gotStrings))
	}

	for i, gotString := range gotStrings {
		if gotString != wantStrings[i] {
			t.Errorf("Got %v, want %v", gotString, wantStrings[i])
		}
	}
}
