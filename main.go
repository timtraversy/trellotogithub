package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/timtraversy/trellotogithub/githubdeviceauth"
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
		githubAuthenticator := githubdeviceauth.NewAuthenticator()
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
