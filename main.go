package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
	"github.com/timtraversy/trellotogithub/deviceauth"
	"golang.org/x/oauth2"
)

type clientFactory struct {
	newTrelloClient        func(token string) trelloClient
	newGithubAuthenticator func() githubAuthenticator
	newGithubClient        func(token string) githubClient
}

type trelloClient interface {
	GetMember(memberID string, args trello.Arguments) (member *trello.Member, err error)
	GetMyBoards(args trello.Arguments) (boards []*trello.Board, err error)
}

type githubAuthenticator interface {
	AuthenticateGithub(scopes []github.Scope) (token string, err error)
}

type githubClient struct {
	repositories repositoriesService
	users        userService
}

type repositoriesService interface {
	List(ctx context.Context, user string, opts *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error)
	ListProjects(ctx context.Context, owner, repo string, opts *github.ProjectListOptions) ([]*github.Project, *github.Response, error)
}

type userService interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
}

func main() {
	in := os.Stdin
	out := os.Stdout
	cliFactory := clientFactory{
		newTrelloClient: func(token string) trelloClient {
			return trello.NewClient(apikey, token)
		},
		newGithubAuthenticator: func() githubAuthenticator {
			return deviceauth.NewAuthenticator()
		},
		newGithubClient: func(token string) githubClient {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			client := github.NewClient(tc)
			return githubClient{
				users:        client.Users,
				repositories: client.Repositories,
			}
		}}
	configure(in, out, cliFactory)
	// TODO: execute import and export
}

type Configuration struct {
	TrelloAuth AuthConfiguration    `yaml:"trello"`
	GithubAuth AuthConfiguration    `yaml:"github"`
	Mapping    MappingConfiguration `yaml:"mapping"`
}

type AuthConfiguration struct {
	Token    string `yaml:"token"`
	Username string `yaml:"username"`
	// Username yaml.Node `yaml:"username"`
}

type MappingConfiguration struct {
	TrelloBoard          string `yaml:"trello_board"`
	TrelloBoardName      string `yaml:"trello_board_name"`
	GithubRepository     string `yaml:"github_repository"`
	GithubRepositoryName string `yaml:"github_repository_name"`
	GithubProject        string `yaml:"github_project"`
	GithubProjectName    string `yaml:"github_project_name"`
}

func configure(in io.Reader, out io.Writer, cliFactory clientFactory) *Configuration {
	fmt.Fprintln(out, "# Configuration")

	// d, err := ioutil.ReadFile("config.yaml")
	// var config Configuration
	// yaml.Unmarshal(d, &config)

	// if err != nil {
	// 	fmt.Fprintln(out, "No configuration file found in this directory. Follow the prompts to configure the tool.")
	// }

	fmt.Fprintln(out, "No configuration file found in this directory. Follow the prompts to configure the tool.")
	fmt.Fprintln(out, "")

	trelloClient, trelloAuth := authenticateTrello(in, out, cliFactory)

	githubClient, githubAuth := authenticateGithub(in, out, cliFactory)

	fmt.Fprintln(out, "## Card mapping")

	board := selectBoard(in, out, trelloClient)

	repository := selectRepository(in, out, githubAuth.Username, githubClient.repositories)

	project := selectProject(in, out, *repository, githubClient.repositories)

	config := Configuration{
		TrelloAuth: trelloAuth,
		GithubAuth: githubAuth,
		Mapping: MappingConfiguration{
			TrelloBoard:          board.ID,
			TrelloBoardName:      board.Name,
			GithubRepository:     fmt.Sprint(repository.GetID()),
			GithubRepositoryName: repository.GetName(),
			GithubProject:        fmt.Sprint(project.GetID()),
			GithubProjectName:    project.GetName(),
		},
	}

	fmt.Fprintln(out, "## Confirm")
	fmt.Fprintln(out, "Your configuration is:")
	fmt.Fprintf(out, "Trello account: %v\n", config.TrelloAuth.Username)
	fmt.Fprintf(out, "Trello board: %v (ID: %v)\n", config.Mapping.TrelloBoardName, config.Mapping.TrelloBoard)
	fmt.Fprintf(out, "GitHub account: %v\n", config.GithubAuth.Username)
	fmt.Fprintf(out, "GitHub repository: %v (ID: %v)\n", config.Mapping.GithubRepositoryName, config.Mapping.GithubRepository)
	fmt.Fprintf(out, "GitHub project: %v (ID: %v)\n\n", config.Mapping.GithubProjectName, config.Mapping.GithubProject)

	fmt.Fprintln(out, "Would you like to save this configuration to a 'config.yaml' file for next time?")
	fmt.Fprint(out, "(y/n) ")

	var yesNo string
	fmt.Fscan(in, yesNo)
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "")

	fmt.Fprintln(out, "Would you like to proceed with the import using these settings?")
	fmt.Fprint(out, "(y/n) ")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "")

	fmt.Fprintln(out, "Configuration complete.")

	return nil
}

func authenticateTrello(in io.Reader, out io.Writer, cliFactory clientFactory) (trelloClient, AuthConfiguration) {
	fmt.Fprintln(out, "## Trello authentication")
	fmt.Fprintf(out, "In your browser, open https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%v\n", apikey)
	fmt.Fprintln(out, "When the authentication process is complete, enter the token you receied here.")
	fmt.Fprint(out, "Token: ")

	var token string
	fmt.Fscan(in, &token)

	client := cliFactory.newTrelloClient(token)
	member, _ := client.GetMember("me", trello.Defaults())

	fmt.Fprintf(out, "Successfully authenticated with Trello as %v.\n\n", member.Username)

	return client, AuthConfiguration{
		Token:    token,
		Username: member.Username,
	}
}

func authenticateGithub(in io.Reader, out io.Writer, cliFactory clientFactory) (githubClient, AuthConfiguration) {
	fmt.Fprintln(out, "## GitHub authentication")
	githubAuthenticator := cliFactory.newGithubAuthenticator()
	githubScopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
	token, _ := githubAuthenticator.AuthenticateGithub(githubScopes)
	// token := "oauth"

	client := cliFactory.newGithubClient(token)

	user, _, _ := client.users.Get(context.Background(), "")

	fmt.Fprintf(out, "Successfully authenticated with GitHub as %v.\n\n", user.GetLogin())

	return client, AuthConfiguration{
		Token:    token,
		Username: user.GetLogin(),
	}
}

func selectBoard(in io.Reader, out io.Writer, client trelloClient) *trello.Board {
	fmt.Fprintln(out, "### Trello board")
	fmt.Fprintln(out, "Fetching Trello boards...")
	boards, _ := client.GetMyBoards(trello.Defaults())
	fmt.Fprintln(out, "Select the Trello board you want to import from:")
	for i, board := range boards {
		fmt.Fprintf(out, "[%v] %v (ID: %v)\n", i, board.Name, board.ID)
	}
	fmt.Fprint(out, "Board: ")

	var selection int
	fmt.Fscan(in, &selection)

	board := boards[selection]
	fmt.Fprintf(out, "Selected '%v'.\n\n", board.Name)

	return boards[selection]
}

func selectRepository(in io.Reader, out io.Writer, user string, client repositoriesService) *github.Repository {
	fmt.Fprintln(out, "### GitHub repository")
	fmt.Fprintln(out, "Fetching GitHub repositories...")
	repos, _, _ := client.List(context.Background(), user, nil)
	fmt.Fprintln(out, "Select the GitHub repository whose projects you want to export to:")
	for i, repo := range repos {
		fmt.Fprintf(out, "[%v] %v (ID: %v)\n", i, repo.GetName(), repo.GetID())
	}
	fmt.Fprint(out, "Repository: ")

	var selection int
	fmt.Fscan(in, &selection)

	repo := repos[selection]
	fmt.Fprintf(out, "Selected '%v'.\n\n", repo.GetName())

	return repos[selection]
}

func selectProject(in io.Reader, out io.Writer, repository github.Repository, client repositoriesService) *github.Project {
	fmt.Fprintln(out, "### GitHub project")
	fmt.Fprintln(out, "Fetching GitHub projects...")
	projects, _, _ := client.ListProjects(context.Background(), repository.GetOwner().GetLogin(), repository.GetName(), nil)
	fmt.Fprintln(out, "Select the GitHub project you want to export to:")
	for i, project := range projects {
		fmt.Fprintf(out, "[%v] %v (ID: %v)\n", i, project.GetName(), project.GetID())
	}
	fmt.Fprint(out, "Project: ")

	var selection int
	fmt.Fscan(in, &selection)

	project := projects[selection]
	fmt.Fprintf(out, "Selected '%v'.\n\n", project.GetName())

	return projects[selection]
}
