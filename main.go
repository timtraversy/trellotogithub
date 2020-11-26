package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/timtraversy/trellotogithub/githubdeviceauth"
)

var (
	// commmands
	authCmd    = flag.NewFlagSet("-auth", flag.ExitOnError)
	authTrello = authCmd.Bool("-trello", false, "Pass this to reauthenticate Trello")
	authGithub = authCmd.Bool("-github", false, "Pass this to reauthenticate GitHub")
)

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
	if len(os.Args) == 1 {
		fmt.Println("Run")
		return
	}
	switch os.Args[1] {
	case authCmd.Name():
		authCmd.Parse(os.Args[2:])
		if *authTrello {
			authenticateTrello(in, out)
		} else if *authGithub {
			authenticateGithub(out)
		} else {
			authenticate(out)
		}
	}

}

func authenticate(out io.Writer) (trelloToken, githubToken string, err error) {
	trelloToken, err = authenticateTrello(os.Stdin, out)
	if err != nil {
		return "", "", err
	}
	githubToken, err = authenticateGithub(out)
	if err != nil {
		return "", "", err
	}
	return trelloToken, githubToken, nil
}

func authenticateGithub(out io.Writer) (token string, err error) {
	githubScopes := []github.Scope{github.ScopeRepo, github.ScopeAdminOrg, github.ScopeUser}
	githubAuthenticator := githubdeviceauth.NewAuthenticator()
	token, err = githubAuthenticator.AuthenticateGithub(githubScopes)
	return
}
