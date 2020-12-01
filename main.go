package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

type TrelloClient interface {
}

type GithubClient interface {
}

func main() {
	in := os.Stdin
	out := os.Stdout
	clientFactory := ClientFactory{
		trelloClientFactory: func(token string) TrelloClient {
			return trello.NewClient(apikey, token)
		},
		githubClientFactory: func(token string) GithubClient {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			return github.NewClient(tc)
		}}
	configure(in, out, clientFactory)
	// TODO: execute import and export
}

type ClientFactory struct {
	trelloClientFactory func(token string) TrelloClient
	githubClientFactory func(token string) GithubClient
}

type Configuration struct {
	TrelloAuth AuthConfiguration    `yaml:"trello"`
	GithubAuth AuthConfiguration    `yaml:"github"`
	Mapping    MappingConfiguration `yaml:"mapping"`
}

type AuthConfiguration struct {
	Token    string    `yaml:"token"`
	Username yaml.Node `yaml:"username"`
}

type MappingConfiguration struct {
	TrelloBoard          string `yaml:"trello_board"`
	TrelloBoardName      string `yaml:"trello_board_name"`
	GithubRepository     string `yaml:"github_repository"`
	GithubRepositoryName string `yaml:"github_repository_name"`
	GithubProject        string `yaml:"github_project"`
	GithubProjectName    string `yaml:"github_project_name"`
}

func configure(in io.Reader, out io.Writer, clientFactory ClientFactory) *Configuration {
	var configuration Configuration
	d, err := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(d, &configuration)

	if err != nil {
		// no config
	}

	// Check missing values
	// if configuration

	// if (configuration == Configuration{}) {
	// 	// S
	// 	// have to auth
	// 	trelloToken, _ := authenticateTrello(in, out)
	// 	githubAuthenticator := deviceauth.NewAuthenticator()
	// 	githubScopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
	// 	githubToken, _ := githubAuthenticator.AuthenticateGithub(githubScopes)
	// 	authTokens = AuthTokens{
	// 		Trello: trelloToken,
	// 		Github: githubToken,
	// 	}
	// 	d, _ = yaml.Marshal(authTokens)
	// 	fmt.Println(string(d))
	// 	ioutil.WriteFile(".auth.yaml", d, 0644)
	// }

	// trelloClient := trello.NewClient(apikey, authTokens.Trello)

	// board := selectBoard(trelloClient, in, out)
	// repo := selectRepository(githubClient.Repositories, in, out)

	// // fmt.Println("Importing from ", board.Name, " to ", project.Name)
	// fmt.Println("Importing from x to", repo.GetName())
	return nil
}

func authenticateTrello(in io.Reader, out io.Writer) (token string, err error) {
	fmt.Fprintln(out, "Open the following URL in your browser to authorize this tool to read the boards on your Trello account.")
	fmt.Fprintf(out, "https://trello.com/1/authorize?expiration=never&name=Trello%%20to%%20Github&scope=read,account&response_type=token&key=%v\n", apikey)
	fmt.Fprintln(out, "When the authentication process is complete, enter the token you receied here:")

	tokenChn := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	tokenChn <- scanner.Text()
	return <-tokenChn, nil
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
