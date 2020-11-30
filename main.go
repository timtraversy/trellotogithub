package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/timtraversy/trellotogithub/deviceauth"

	"github.com/adlio/trello"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

// var (
// 	// commmands
// 	authCmd    = flag.NewFlagSet("-auth", flag.ExitOnError)
// 	authTrello = authCmd.Bool("-trello", false, "Pass this to reauthenticate Trello")
// 	authGithub = authCmd.Bool("-github", false, "Pass this to reauthenticate GitHub")
// )

var (
	exitCode = 0
)

func main() {
	deviceauth.NewAuthenticator()
	return
	githubToTrelloMain()
	os.Exit(exitCode)
}

func githubToTrelloMain() {
	in := os.Stdin
	out := os.Stdout
	// if len(os.Args) == 1 {
	// 	run()
	// }
	// switch os.Args[1] {
	// case authCmd.Name():
	// 	authCmd.Parse(os.Args[2:])
	// 	if *authTrello {
	// 		authenticateTrello(in, out)
	// 	} else if *authGithub {
	// 		authenticateGithub(out)
	// 	} else {
	// 		authenticate(out)
	// 	}
	// }
	var authTokens AuthTokens
	d, _ := ioutil.ReadFile(".auth.yaml")
	yaml.Unmarshal(d, &authTokens)
	if (authTokens == AuthTokens{}) {
		// have to auth
		trelloToken, _ := authenticateTrello(in, out)
		githubAuthenticator := deviceauth.NewAuthenticator()
		githubScopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
		githubToken, _ := githubAuthenticator.AuthenticateGithub(githubScopes)
		authTokens = AuthTokens{
			Trello: trelloToken,
			Github: githubToken,
		}
		d, _ = yaml.Marshal(authTokens)
		fmt.Println(string(d))
		ioutil.WriteFile(".auth.yaml", d, 0644)
	}

	// trelloClient := trello.NewClient(apikey, authTokens.Trello)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authTokens.Github},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	githubClient := github.NewClient(tc)

	// board := selectBoard(trelloClient, in, out)
	repo := selectRepository(githubClient.Repositories, in, out)

	// fmt.Println("Importing from ", board.Name, " to ", project.Name)
	fmt.Println("Importing from x to", repo.GetName())
}

type AuthTokens struct {
	Trello string `yaml:"trello"`
	Github string `yaml:"github"`
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

// func authenticate(out io.Writer) (trelloToken, githubToken string, err error) {
// 	trelloToken, err = authenticateTrello(os.Stdin, out)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	githubToken, err = authenticateGithub(out)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return trelloToken, githubToken, nil
// }

// func authenticateGithub(out io.Writer) (token string, err error) {
// 	githubScopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
// 	githubAuthenticator := githubdeviceauth.NewAuthenticator()
// 	token, err = githubAuthenticator.AuthenticateGithub(githubScopes)
// 	return
// }
